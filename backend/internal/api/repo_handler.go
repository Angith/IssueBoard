package api

import (
	"encoding/json"
	"net/http"

	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/angith/issueboard/internal/service"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
)

type RepoHandler struct {
	repoService *service.RepoService
}

func NewRepoHandler(repoService *service.RepoService) *RepoHandler {
	return &RepoHandler{repoService: repoService}
}

func (h *RepoHandler) AddRepository(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*types.User)
	userID, _ := uuid.Parse(user.ID.String())

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
	user := r.Context().Value(middleware.UserKey).(*types.User)
	userID, _ := uuid.Parse(user.ID.String())

	repos, err := h.repoService.ListRepositories(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}
