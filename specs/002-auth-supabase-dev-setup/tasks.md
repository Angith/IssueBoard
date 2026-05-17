# Tasks: Supabase Authentication and Dev Experience Refactor

**Input**: Design documents from `/specs/002-auth-supabase-dev-setup/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Unit tests for Go backend and React frontend are included as per Constitution Requirements IV.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 [P] Configure backend environment variable loader in `backend/internal/config/config.go`
- [X] T002 [P] Initialize Supabase client in frontend at `frontend/src/lib/supabase.ts`
- [X] T003 [P] Add Docker Compose for local PostgreSQL in `infra/docker-compose.yml`
- [X] T004 [P] Update `backend/go.mod` with `golang-jwt/jwt`, `joho/godotenv`, and `caarlos0/env`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T005 Create many-to-many database migration for `user_repository` join table in `backend/migrations/002_user_repo_many_to_many.sql`
- [X] T006 Update backend User and Repository models to align with Supabase sub IDs and shared repos in `backend/internal/repository/models/user.go` and `backend/internal/repository/models/repository.go`
- [X] T007 [P] Implement Supabase JWT validation middleware in `backend/internal/api/middleware/auth.go`
- [X] T008 [P] Implement central config loader in `backend/cmd/main.go` using `.env`
- [X] T009 [P] Create GitHub client wrapper for unauthenticated requests in `backend/internal/github/client.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Passwordless Login Journey (Priority: P1) 🎯 MVP

**Goal**: Implement email-based magic link login using Supabase.

**Independent Test**: User can enter email on `/login`, receive a link (visible in local Supabase dashboard), and become authenticated.

### Tests for User Story 1

- [X] T010 [P] [US1] Unit test for JWT validation middleware in `backend/internal/api/middleware/auth_test.go`
- [X] T011 [US1] Integration test for authentication state management in `frontend/src/app/login/login.test.tsx`

### Implementation for User Story 1

- [X] T012 [P] [US1] Create login page with email input in `frontend/src/app/login/page.tsx`
- [X] T013 [P] [US1] Implement auth state provider and protected routes in `frontend/src/components/AuthProvider.tsx`
- [X] T014 [US1] Replace GitHub OAuth routes with Supabase session handling in `backend/internal/api/issue_handler.go` and `backend/internal/api/repo_handler.go`
- [X] T015 [P] [US1] Implement logout functionality in `frontend/src/components/Navbar.tsx`
- [X] T016 [US1] Update API client to attach Authorization header in `frontend/src/services/api.ts`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 2 - User-Scoped Repository Inventory (Priority: P1)

**Goal**: Manage personal repository lists with data isolation using many-to-many logic.

**Independent Test**: User A adds Repo X; User B logs in and sees an empty list. User B adds Repo X; Repo X is linked to User B without duplication.

### Tests for User Story 2

- [X] T017 [US1] Unit test for many-to-many repository linking in `backend/internal/repository/repo_repository_test.go`
- [X] T018 [US1] Integration test for user-scoped repo listing in `backend/tests/integration/repo_test.go`

### Implementation for User Story 2

- [X] T019 [P] [US1] Update Repository repository to support many-to-many linking in `backend/internal/repository/repo_repository.go`
- [X] T020 [US1] Refactor Repository service to inject `user_id` from context in `backend/internal/service/repo_service.go`
- [X] T021 [US1] Implement full URL parsing for repository addition in `backend/internal/api/repo_handler.go`
- [X] T022 [P] [US1] Update repository list component to show only user's repos in `frontend/src/app/inventory/page.tsx`

**Checkpoint**: User Stories 1 AND 2 should both work independently.

---

## Phase 5: User Story 3 - Streamlined Local Setup (Priority: P1)

**Goal**: Enable 15-minute project setup for new developers.

**Independent Test**: Following README instructions from scratch results in a running app with local Auth.

### Implementation for User Story 3

- [X] T023 [P] [US3] Create `README.md` with copy-paste setup commands and architecture update
- [X] T024 [P] [US3] Create `quickstart.md` with detailed local Supabase CLI instructions in `specs/002-auth-supabase-dev-setup/quickstart.md`
- [X] T025 [US3] Implement rate-limit graceful handling and feedback in `frontend/src/components/IssueBoard.tsx`
- [X] T026 [P] [US3] Add "Refresh" button for on-demand issue updates in `frontend/src/app/repos/[id]/page.tsx`

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Cleanup and final validation

- [X] T027 [P] Remove all legacy GitHub OAuth code from `backend/internal/api/middleware/auth_old.go` (and related files)
- [X] T028 [P] Remove `oauth_token` and `github_id` columns from `users` table in `backend/migrations/003_cleanup_legacy_auth.sql`
- [X] T029 Code cleanup and refactoring in `backend/internal/github/` to ensure unauthenticated-only logic
- [X] T030 Final run of `specs/002-auth-supabase-dev-setup/quickstart.md` validation
- [X] T031 [US3] Measure and verify login-to-dashboard flow duration (SC-002)
- [X] T032 [P] Exhaustive security audit for user data isolation (SC-004)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Foundation ready - No story dependencies
- **User Story 2 (P1)**: Foundation ready - Depends on Auth context from US1 for full testing, but backend can be implemented independently
- **User Story 3 (P1)**: Foundation ready - Can be implemented in parallel with others

---

## Implementation Strategy

### MVP First (User Story 1 & 2)

1. Complete Phase 1 & 2 (Foundation)
2. Complete Phase 3 (Auth)
3. Complete Phase 4 (Scoped Repos)
4. **STOP and VALIDATE**: Verify end-to-end flow from login to repo management.

### Parallel Opportunities

- T001, T002, T003, T004 can run in parallel.
- Once Foundation (T005-T009) is complete, US1, US2, and US3 implementation can proceed in parallel.
- UI components (T012, T013, T015) can be built while backend handlers (T014) are being refactored.
