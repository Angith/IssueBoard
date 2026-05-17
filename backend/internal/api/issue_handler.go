package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

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
		logrus.WithError(err).WithField("repo_id", repoID).Error("Failed to get categorized issues")
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

	logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID}).Debug("Refreshing issues from GitHub")
	if err := h.issueService.RefreshIssues(r.Context(), userID, repoID); err != nil {
		logrus.WithError(err).WithField("repo_id", repoID).Error("Failed to refresh issues")
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

	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Getting available labels")

	labels, err := h.issueService.GetAvailableLabels(r.Context(), userID, repoID)
	if err != nil {
		log.WithError(err).Error("Failed to get available labels")
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

	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID})
	log.Debug("Getting tracked labels")

	labels, err := h.issueService.GetTrackedLabels(r.Context(), userID, repoID)
	if err != nil {
		log.WithError(err).Error("Failed to get tracked labels")
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
		logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID}).WithError(err).Warn("Invalid request body for UpdateTrackedLabels")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log := logrus.WithFields(logrus.Fields{"user_id": userID, "repo_id": repoID, "label_count": len(body.Labels)})
	log.Debug("Updating tracked labels")

	if err := h.issueService.UpdateTrackedLabels(r.Context(), userID, repoID, body.Labels); err != nil {
		log.WithError(err).Error("Failed to update tracked labels")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Tracked labels updated")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
