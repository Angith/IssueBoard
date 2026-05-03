package api

import (
	"encoding/json"
	"net/http"

	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/angith/issueboard/internal/service"
	"github.com/google/uuid"
)

type RepoHandler struct {
	repoService *service.RepoService
}

func NewRepoHandler(repoService *service.RepoService) *RepoHandler {
	return &RepoHandler{repoService: repoService}
}

func (h *RepoHandler) AddRepository(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	repo, err := h.repoService.AddRepository(r.Context(), userID, body.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
}

func (h *RepoHandler) ListRepositories(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	repos, err := h.repoService.ListRepositories(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}
