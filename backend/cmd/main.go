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
	"github.com/supabase-community/gotrue-go"
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

	// Initialize services
	// We need a GitHub client with a token, usually passed in headers or fetched from DB
	// For simplicity in main, we'll initialize the client per-request in the handlers
	// or use a service that takes a client.

	// Mocking GitHub client init
	ghClient := github.NewClient("fake-token")
	ghService := github.NewRepoService(ghClient)

	repoService := service.NewRepoService(repoRepo, ghService)
	issueService := service.NewIssueService(issueRepo, repoRepo, ghService)

	// Initialize handlers
	repoHandler := api.NewRepoHandler(repoService)
	issueHandler := api.NewIssueHandler(issueService)

	// Initialize Supabase Auth client for middleware
	authClient := gotrue.New(cfg.SupabaseURL, cfg.SupabaseAnonKey)

	// Set up routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Protected routes
	authMidd := middleware.AuthMiddleware(authClient)

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
