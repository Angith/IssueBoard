package api

import (
	"encoding/json"
	"net/http"

	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/angith/issueboard/internal/crypto"
	"github.com/angith/issueboard/internal/repository"
	"github.com/google/go-github/v60/github"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userRepo      *repository.UserRepository
	encryptionKey string
}

func NewUserHandler(userRepo *repository.UserRepository, encryptionKey string) *UserHandler {
	return &UserHandler{
		userRepo:      userRepo,
		encryptionKey: encryptionKey,
	}
}

type SetGitHubTokenRequest struct {
	Token string `json:"token"`
}

func (h *UserHandler) SetGitHubToken(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SetGitHubTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Validate token with GitHub
	client := github.NewClient(nil).WithAuthToken(req.Token)
	_, _, err := client.Users.Get(r.Context(), "")
	if err != nil {
		logrus.WithError(err).Warn("Failed to validate GitHub token")
		http.Error(w, "Invalid GitHub token", http.StatusUnauthorized)
		return
	}

	// Encrypt
	encrypted, err := crypto.Encrypt(req.Token, h.encryptionKey)
	if err != nil {
		logrus.WithError(err).Error("Failed to encrypt token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.userRepo.SaveGitHubToken(r.Context(), userID, encrypted); err != nil {
		logrus.WithError(err).Error("Failed to save github token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *UserHandler) DeleteGitHubToken(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.userRepo.DeleteGitHubToken(r.Context(), userID); err != nil {
		logrus.WithError(err).Error("Failed to delete github token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *UserHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hasToken := false
	token, err := h.userRepo.GetGitHubToken(r.Context(), userID)
	if err == nil && len(token) > 0 {
		hasToken = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"has_github_token": hasToken,
	})
}
