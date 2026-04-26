# Tasks: Issue Board Core

**Input**: Design documents from `/specs/001-issue-board-core/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Backend unit tests are REQUIRED (70% coverage goal per constitution and spec).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/src/`, `frontend/src/`
- Paths below follow the monorepo structure defined in plan.md

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 [P] Create backend project structure in `backend/`
- [X] T002 [P] Create frontend project structure in `frontend/`
- [X] T003 [P] Configure backend linting and formatting in `backend/.golangci.yml`
- [X] T004 [P] Configure frontend linting and formatting in `frontend/.eslintrc.json` and `frontend/.prettierrc`
- [X] T005 [P] Setup `infra/docker-compose.yml` for local PostgreSQL

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- [X] T006 Implement environment configuration system in `backend/internal/config/config.go`
- [X] T007 Implement database connection layer with `pgx` in `backend/internal/repository/db.go`
- [X] T008 Setup database migration framework in `backend/migrations/`
- [X] T009 [P] Initialize GitHub API client abstraction in `backend/internal/github/client.go` using `google/go-github`
- [X] T010 [P] Setup Next.js project with Supabase Auth client in `frontend/src/lib/supabase.ts`

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - GitHub Authentication (Priority: P1) 🎯 MVP

**Goal**: Enable users to log in using their GitHub account via Supabase.

**Independent Test**: Successfully redirect to GitHub, authorize, and return to a secure session in the app.

### Implementation for User Story 1

- [X] T011 [P] [US1] Create User model in `backend/internal/repository/models/user.go`
- [X] T012 [US1] Implement User repository in `backend/internal/repository/user_repository.go`
- [X] T013 [US1] Implement Authentication service in `backend/internal/service/auth_service.go`
- [X] T014 [US1] Create Login page in `frontend/src/pages/login.tsx` using Supabase Auth UI
- [X] T015 [US1] Implement Auth middleware in `backend/internal/api/middleware/auth.go`
- [X] T016 [US1] Unit test for Authentication service in `backend/internal/service/auth_service_test.go` (mocking Supabase/DB)

**Checkpoint**: User Story 1 complete - Authentication is functional.

---

## Phase 4: User Story 2 - Repository Inventory Management (Priority: P1)

**Goal**: Allow users to add and list GitHub repositories in their private inventory.

**Independent Test**: Add a valid GitHub URL and see it appear in the inventory list.

### Implementation for User Story 2

- [X] T017 [P] [US2] Create Repository model in `backend/internal/repository/models/repository.go`
- [X] T018 [US2] Implement Repository repository in `backend/internal/repository/repo_repository.go`
- [X] T019 [US2] Implement GitHub repository validation in `backend/internal/github/repo_service.go`
- [X] T020 [US2] Implement Repository service in `backend/internal/service/repo_service.go`
- [X] T021 [US2] Implement `POST /api/repos` and `GET /api/repos` handlers in `backend/internal/api/repo_handler.go`
- [X] T022 [US2] Create Inventory page in `frontend/src/pages/inventory.tsx`
- [X] T023 [US2] Create Add Repository component in `frontend/src/components/AddRepoForm.tsx`
- [X] T024 [US2] Unit test for Repository service in `backend/internal/service/repo_service_test.go` (mocking GitHub API and DB)

**Checkpoint**: User Story 2 complete - Inventory management is functional.

---

## Phase 5: User Story 3 - Categorized Issue Board View (Priority: P1)

**Goal**: Display open issues for a repository grouped by labels.

**Independent Test**: Select a repo and see issues categorized by labels (including "Unlabeled").

### Implementation for User Story 3

- [X] T025 [P] [US3] Create Issue and Label models in `backend/internal/repository/models/issue.go` and `label.go`
- [X] T026 [US3] Implement Issue and Label repositories in `backend/internal/repository/issue_repository.go`
- [X] T027 [US3] Implement Issue grouping logic (including duplication for multi-label support) in `backend/internal/service/issue_service.go`
- [X] T028 [US3] Implement `GET /api/repos/{id}/issues` handler in `backend/internal/api/issue_handler.go`
- [X] T029 [US3] Create Issue Board page in `frontend/src/pages/repos/[id].tsx`
- [X] T030 [US3] Create Issue Board component in `frontend/src/components/IssueBoard.tsx`
- [X] T031 [US3] Unit test for Issue service in `backend/internal/service/issue_service_test.go` (mocking DB and GitHub)

**Checkpoint**: User Story 3 complete - Categorized board view is functional.

---

## Phase 6: User Story 4 - Issue Caching & On-Demand Refresh (Priority: P2)

**Goal**: Use cached data for performance and provide a manual refresh option.

**Independent Test**: Load board instantly from cache; click refresh to trigger GitHub API update.

### Implementation for User Story 4

- [X] T032 [US4] Implement Cache-first fetch logic in `backend/internal/service/issue_service.go`
- [X] T033 [US4] Implement `POST /api/repos/{id}/refresh` handler in `backend/internal/api/issue_handler.go`
- [X] T034 [US4] Add Refresh button to `frontend/src/components/IssueBoard.tsx`
- [X] T035 [US4] Unit test for Cache-first logic in `backend/internal/service/issue_service_test.go`

**Checkpoint**: User Story 4 complete - Performance caching and manual refresh are functional.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final refinements and quality checks.

- [X] T036 [P] Ensure all backend services have > 70% test coverage
- [X] T037 [P] Final linting check across backend and frontend
- [X] T038 Add loading and error states to all frontend pages
- [X] T039 Verify environment-based DB configuration works for both local and production
- [X] T040 Complete `quickstart.md` validation by running through setup steps
- [X] T041 Verify performance success criteria (SC-001, SC-002) under simulated load

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Can start immediately.
- **Foundational (Phase 2)**: Depends on Phase 1 completion.
- **User Stories (Phase 3-6)**: All depend on Phase 2 completion.
  - US1 (Auth) should be completed first as US2 depends on user context.
  - US2 (Inventory) must be completed before US3 (Board).
  - US3 (Board) must be completed before US4 (Caching).

### Parallel Opportunities

- T001-T005 in Phase 1 can run in parallel.
- T009 and T010 in Phase 2 can run in parallel.
- T011 and T014 in Phase 3 can start in parallel.
- T017, T022, T023 in Phase 4 can start in parallel.
- T025, T029, T030 in Phase 5 can start in parallel.

---

## Implementation Strategy

### MVP First (User Stories 1-3)

1. Complete Setup and Foundational phases.
2. Complete US1 (Auth), US2 (Inventory), and US3 (Board) to achieve core value.
3. **VALIDATE**: Ensure a user can log in, add a repo, and see issues grouped by labels.

### Incremental Delivery

- Each user story is an independent increment that adds value.
- Caching (US4) is treated as a P2 performance enhancement after core functionality is verified.
