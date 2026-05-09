package api

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("GetCategorizedIssues error: %v\n", err)
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

func (h *IssueHandler) GetAvailableLabels(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// URL format: /api/repos/{id}/labels/available
	idStr := r.URL.Path[len("/api/repos/") : len(r.URL.Path)-len("/labels/available")]
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	labels, err := h.issueService.GetAvailableLabels(r.Context(), userID, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(labels)
}

func (h *IssueHandler) GetTrackedLabels(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// URL format: /api/repos/{id}/labels/tracked
	idStr := r.URL.Path[len("/api/repos/") : len(r.URL.Path)-len("/labels/tracked")]
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	labels, err := h.issueService.GetTrackedLabels(r.Context(), userID, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(labels)
}

func (h *IssueHandler) UpdateTrackedLabels(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// URL format: /api/repos/{id}/labels/tracked
	idStr := r.URL.Path[len("/api/repos/") : len(r.URL.Path)-len("/labels/tracked")]
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Labels []string `json:"labels"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.issueService.UpdateTrackedLabels(r.Context(), userID, repoID, body.Labels); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
