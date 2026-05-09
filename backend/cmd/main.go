package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/angith/issueboard/internal/api"
	"github.com/angith/issueboard/internal/api/middleware"
	"github.com/angith/issueboard/internal/config"
	"github.com/angith/issueboard/internal/github"
	"github.com/angith/issueboard/internal/repository"
	"github.com/angith/issueboard/internal/service"
)

func main() {
	cfg := config.Load()

	dbPool, err := repository.NewDBPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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

	// Initialize handlers
	repoHandler := api.NewRepoHandler(repoService)
	issueHandler := api.NewIssueHandler(issueService)

	// Set up routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Protected routes
	authMidd := middleware.AuthMiddleware(cfg.SupabaseURL, userRepo)

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

	// Wrap mux with CORS middleware
	handler := middleware.CorsMiddleware(mux)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
