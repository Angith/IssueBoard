package api

import (
	"encoding/json"
	"net/http"

	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/angith/issueboard/internal/service"
	"github.com/google/uuid"
)

type IssueHandler struct {
	issueService *service.IssueService
}

func NewIssueHandler(issueService *service.IssueService) *IssueHandler {
	return &IssueHandler{issueService: issueService}
}

func (h *IssueHandler) GetCategorizedIssues(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/api/repos/") : len(r.URL.Path)-len("/issues")]
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	board, err := h.issueService.GetCategorizedIssues(r.Context(), userID, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func (h *IssueHandler) RefreshIssues(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/api/repos/") : len(r.URL.Path)-len("/refresh")]
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	if err := h.issueService.RefreshIssues(r.Context(), userID, repoID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
