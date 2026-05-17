# Implementation Plan: Supabase Authentication and Dev Experience Refactor

**Branch**: `002-auth-supabase-dev-setup` | **Date**: 2026-05-02 | **Spec**: [specs/002-auth-supabase-dev-setup/spec.md](spec.md)
**Input**: Feature specification for refactoring authentication and local developer experience.

## Summary
Replace the existing GitHub OAuth system with Supabase Magic Link authentication and streamline the local development environment using Docker and Supabase CLI. The refactor focuses on user-scoped data isolation and removing all dependencies on GitHub user tokens.

## Technical Context

**Language/Version**: Go 1.25.0 (Backend), React 19 / Next.js 16 (Frontend)
**Primary Dependencies**: `supabase-js` (Frontend), `gotrue-go` (Backend), `pgx/v5` (Backend), `google/go-github/v60` (GitHub API)
**Storage**: PostgreSQL (via Supabase)
**Testing**: `go test` (Backend), standard React testing tools (Frontend)
**Target Platform**: Web
**Project Type**: Full-stack web application
**Performance Goals**: <60s login-to-dashboard flow
**Constraints**: Purely unauthenticated GitHub API (60 req/hr limit), Docker requirement for local dev
**Scale/Scope**: Data isolation for multiple users, shared repository records

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Requirement I (Architecture)**: Backend follows Layered Architecture (handlers, services, repositories). **PASS**
- **Requirement II (API-First)**: New auth flow and user-scoped endpoints require contract definitions. **PASS**
- **Requirement III (Clean Code)**: Idiomatic Go and React hooks to be used. **PASS**
- **Requirement IV (Testing)**: 70% coverage goal for new auth and data logic. **PASS**
- **Requirement V (Security)**: Secrets moved to `.env`, Supabase JWT verification on backend. **PASS**

## Project Structure

### Documentation (this feature)

```text
specs/002-auth-supabase-dev-setup/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
backend/
├── cmd/main.go          # Config loader update
├── internal/
│   ├── api/
│   │   ├── middleware/  # New Auth middleware
│   │   └── ...          # Updated user-scoped handlers
│   ├── config/          # .env loader
│   ├── repository/      # Updated many-to-many logic
│   └── service/         # Auth and Repo service refactor
└── migrations/          # New user_repository join table

frontend/
├── src/
│   ├── app/
│   │   ├── login/       # New email login page
│   │   └── ...          # Protected routes logic
│   ├── components/      # Updated Auth guards
│   ├── lib/             # Supabase client init
│   └── services/        # API client with token attachment
└── ...
```

**Structure Decision**: Maintaining the existing `backend/` and `frontend/` separation as per Constitution.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Many-to-Many Repo model | Support shared public repos across users | One-to-many would duplicate global repo data and waste API quota |
