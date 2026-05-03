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
	ghClient := github.NewClient("")
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
	authMidd := middleware.AuthMiddleware(cfg.SupabaseJWTSecret, userRepo)

	mux.Handle("/api/repos", authMidd(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			repoHandler.AddRepository(w, r)
		} else {
			repoHandler.ListRepositories(w, r)
		}
	})))

	mux.Handle("/api/repos/", authMidd(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Very simple router logic
		if r.URL.Path[len(r.URL.Path)-len("/issues"):] == "/issues" {
			issueHandler.GetCategorizedIssues(w, r)
		} else if r.URL.Path[len(r.URL.Path)-len("/refresh"):] == "/refresh" {
			issueHandler.RefreshIssues(w, r)
		}
	})))

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
