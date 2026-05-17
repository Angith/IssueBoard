package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/angith/issueboard/internal/api"
	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/angith/issueboard/internal/config"
	"github.com/angith/issueboard/internal/github"
	"github.com/angith/issueboard/internal/logger"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/service"
)

func main() {
	// Bootstrap with error level so config-load warnings are visible.
	// logger.Init() will re-apply the correct level once config is loaded.
	logger.Init(logger.Error)

	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Error("Failed to load configuration, continuing with minimal defaults")
		cfg = &config.Config{Port: "8080"} // minimal to keep server running
	}

	// Re-initialize with the level from config (LOG_LEVEL env var).
	logger.Init(cfg.LogLevel)

	dbPool, err := repository.NewDBPool(cfg.DatabaseURL)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Initialize repositories
	repoRepo := repository.NewRepoRepository(dbPool)
	issueRepo := repository.NewIssueRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)

	// Initialize GitHub services
	ghClient := github.NewClient(cfg.GitHubToken)
	ghRepoService := github.NewRepoService(ghClient)
	ghIssueService := github.NewIssueService(ghClient)

	// Initialize business services
	repoService := service.NewRepoService(repoRepo, ghRepoService)
	issueService := service.NewIssueService(issueRepo, repoRepo, ghIssueService)

	// Protected routes
	authMidd := middleware.AuthMiddleware(cfg.SupabaseURL, userRepo, cfg.EncryptionKey)

	// Initialize handlers
	repoHandler := api.NewRepoHandler(repoService)
	issueHandler := api.NewIssueHandler(issueService)
	userHandler := api.NewUserHandler(userRepo, cfg.EncryptionKey)

	// Set up routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// User Settings Routes
	mux.Handle("/api/user/settings", authMidd(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetSettings(w, r)
		}
	})))

	mux.Handle("/api/user/settings/github-token", authMidd(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut || r.Method == http.MethodPost {
			userHandler.SetGitHubToken(w, r)
		} else if r.Method == http.MethodDelete {
			userHandler.DeleteGitHubToken(w, r)
		}
	})))

	mux.Handle("/api/repos", authMidd(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			repoHandler.AddRepository(w, r)
		} else {
			repoHandler.ListRepositories(w, r)
		}
	})))

	mux.Handle("/api/repos/", authMidd(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle DELETE /api/repos/{id}
		if r.Method == http.MethodDelete {
			repoHandler.RemoveRepository(w, r)
			return
		}
		// Handle other subpaths
		if len(r.URL.Path) > len("/labels/available") && r.URL.Path[len(r.URL.Path)-len("/labels/available"):] == "/labels/available" {
			issueHandler.GetAvailableLabels(w, r)
		} else if len(r.URL.Path) > len("/labels/tracked") && r.URL.Path[len(r.URL.Path)-len("/labels/tracked"):] == "/labels/tracked" {
			if r.Method == http.MethodPut {
				issueHandler.UpdateTrackedLabels(w, r)
			} else {
				issueHandler.GetTrackedLabels(w, r)
			}
		} else if len(r.URL.Path) > len("/issues") && r.URL.Path[len(r.URL.Path)-len("/issues"):] == "/issues" {
			issueHandler.GetCategorizedIssues(w, r)
		} else if len(r.URL.Path) > len("/refresh") && r.URL.Path[len(r.URL.Path)-len("/refresh"):] == "/refresh" {
			issueHandler.RefreshIssues(w, r)
		}
	})))

	// Wrap mux with middlewares
	handler := middleware.CorsMiddleware(mux)
	handler = middleware.RecoveryMiddleware(handler)
	handler = middleware.LoggerMiddleware(handler)

	logrus.Infof("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		logrus.Fatalf("Server stopped unexpectedly: %v", err)
	}
}
