package api

import (
	"encoding/json"
	"net/http"

	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/angith/issueboard/internal/service"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
		logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).Warn("Invalid request body for AddRepository")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log := logrus.WithFields(logrus.Fields{"user_id": userID, "url": body.URL})
	log.Debug("Adding repository")

	repo, err := h.repoService.AddRepository(r.Context(), userID, body.URL)
	if err != nil {
		log.WithError(err).Error("Failed to add repository")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithField("repo", repo.FullName).Info("Repository added")
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

	log := logrus.WithField("user_id", userID)
	log.Debug("Listing repositories")

	repos, err := h.repoService.ListRepositories(r.Context(), userID)
	if err != nil {
		log.WithError(err).Error("Failed to list repositories")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func (h *RepoHandler) RemoveRepository(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/api/repos/"):]
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Removing repository")

	err = h.repoService.RemoveRepository(r.Context(), userID, repoID)
	if err != nil {
		log.WithError(err).Error("Failed to remove repository")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Repository removed")
	w.WriteHeader(http.StatusNoContent)
}
