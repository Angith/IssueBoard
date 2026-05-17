# Feature Specification: Supabase Authentication and Dev Experience Refactor

**Feature Branch**: `002-auth-supabase-dev-setup`  
**Created**: 2026-05-02  
**Status**: Draft  
**Input**: User description: "Update the existing issueBoard specification to replace GitHub OAuth with a lightweight authentication system using Supabase magic link (email-based login), and improve local development setup and documentation. Context: * Existing system uses GitHub OAuth (to be removed) * Application will now support only public GitHub repositories * Authentication is required only for user identity (not GitHub access) * Backend: Go * Frontend: Next.js * Database: PostgreSQL (Supabase in production, Docker locally) Goal: Refactor authentication and developer experience while keeping the core feature intact (repo inventory + issue board). --- ## 1. Authentication Update (Critical Change) Replace GitHub OAuth with Supabase Auth using magic link (passwordless email login). Requirements: * User logs in via email (magic link) * No passwords required * Supabase handles authentication and session management * Frontend receives session (JWT or access token) * Backend validates user identity via Supabase JWT Backend Responsibilities: * Verify Supabase JWT on each request * Extract user ID from token * Associate all data (repos, issues) with authenticated user Frontend Responsibilities: * Implement login flow (email input → magic link) * Handle session persistence * Attach auth token to API requests Remove: * GitHub OAuth flow * GitHub token storage * Any dependency on GitHub user identity --- ## 2. GitHub Integration Update * Only support public repositories * Use unauthenticated GitHub API requests OR optional server-side token * No user-level GitHub permissions required --- ## 3. Local Development Experience (High Priority) Define a clean and simple local setup: Requirements: * Local PostgreSQL via Docker * Supabase Auth integration must work in local development OR provide a clear fallback/mock mode Environment Configuration: * Use `.env` for all configs: * DATABASE_URL * SUPABASE_URL * SUPABASE_ANON_KEY * SUPABASE_JWT_SECRET (for backend validation) * No hardcoded values --- ## 4. Backend Changes (Go) * Add middleware to validate Supabase JWT * Extract user identity from token * Refactor all handlers to be user-scoped * Ensure all DB queries include user_id filtering Add: * Auth middleware layer * Config loader for environment variables * Clear separation between auth logic and business logic --- ## 5. Frontend Changes (Next.js) * Replace GitHub login with email login UI * Add: * Login page * Auth state management * Protected routes * Ensure token is attached to API requests * Handle: * Loading states * Auth errors * Session persistence --- ## 6. Data Model Updates Ensure user model aligns with Supabase: * user_id (from Supabase) * repositories linked to user_id No GitHub user data should be stored. --- ## 7. README & Documentation (High Priority) Update README to include: 1. Project Overview 2. Tech Stack (Go, React, PostgreSQL, Supabase, Docker) 3. Architecture Summary (updated without GitHub OAuth) 4. Local Development Setup: * How to run PostgreSQL via Docker * How to configure `.env` * How to run backend + frontend * How to enable mock auth (if applicable) 5. Authentication Flow (magic link explanation) 6. Basic Usage: * Login → Add repo → View issues Documentation must be: * Clear * Beginner-friendly * Copy-paste ready commands --- ## 8. Acceptance Criteria * User can log in via email magic link * Backend correctly identifies user via Supabase JWT * User can add repositories and see only their own data * GitHub OAuth is fully removed * Application works in local environment with minimal setup * README enables a new developer to run the project without confusion --- ## 9. Constraints * No GitHub OAuth usage * No private repo support * All authentication must be handled via Supabase or mock mode * Maintain clean, modular architecture --- ## Output Requirements * Clearly structured spec update (not a full rewrite) * Highlight modified sections: * Authentication * Backend * Frontend * Local setup * Documentation * Use precise, implementation-ready language * Ensure compatibility with existing system where possible"

## Clarifications

### Session 2026-05-02
- Q: Local Authentication Strategy → A: Local Supabase CLI
- Q: Repository Identification Method → A: Full URL with auto-parsing
- Q: Duplicate Repository Handling → A: Link existing record to user
- Q: GitHub API Authentication Strategy → A: Purely Unauthenticated
- Q: Data Refresh Strategy → A: On-demand (User Triggered)
- Q: User Profile Information → A: Email only

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Passwordless Login Journey (Priority: P1)

As a user, I want to log into the application using only my email address so that I don't have to remember another password or grant GitHub account permissions.

**Why this priority**: Essential for the new authentication model. It removes friction and dependencies on GitHub OAuth.

**Independent Test**: Can be tested by entering an email on the login page, receiving a magic link, and confirming that clicking the link results in an authenticated session.

**Acceptance Scenarios**:

1. **Given** I am on the login page, **When** I enter a valid email and submit, **Then** I should see a confirmation message that a magic link has been sent.
2. **Given** I have received a magic link email, **When** I click the link, **Then** I should be redirected to my repository inventory, fully authenticated.

---

### User Story 2 - User-Scoped Repository Inventory (Priority: P1)

As an authenticated user, I want to manage my own list of public repositories so that I can track issues relevant to my work.

**Why this priority**: Core functionality of the application. Ensures data isolation between users.

**Independent Test**: Can be tested by adding a public repository and verifying it appears in the list, then logging in as a different user to verify the list is unique to each user.

**Acceptance Scenarios**:

1. **Given** I am logged in, **When** I add a valid public GitHub repository URL, **Then** it should appear in my inventory with its current issue count.
2. **Given** I have added repositories, **When** another user logs in, **Then** they should not see my repositories in their inventory.

---

### User Story 3 - Streamlined Local Setup (Priority: P1)

As a developer, I want to set up the project locally with minimal steps so that I can start contributing or testing features quickly.

**Why this priority**: High priority for developer experience and project maintainability.

**Independent Test**: Can be tested by following the README instructions on a clean machine and verifying the application runs successfully with local Docker services.

**Acceptance Scenarios**:

1. **Given** I have cloned the repo and have Docker installed, **When** I follow the setup guide (e.g., `docker-compose up` and `.env` configuration), **Then** both the backend and frontend should start and connect to the local database.

---

### Edge Cases

- **Expired Magic Link**: If a user clicks an expired link, they should receive a clear error message and an option to resend the link.
- **Invalid/Private Repository**: If a user tries to add a private repository or an invalid URL, the system should return a user-friendly error explaining that only public repositories are supported.
- **Offline Development**: How the system behaves when the Supabase Auth service is unreachable (e.g., during local development without internet).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support passwordless email login (Magic Link) for a frictionless onboarding experience.
- **FR-002**: System MUST securely validate user identity on every backend request using industry-standard secure tokens.
- **FR-003**: System MUST isolate user data, ensuring that the *inventory of tracked repositories* is strictly scoped to the authenticated user.
- **FR-004**: System MUST support tracking of public GitHub repositories via full URL entry (auto-parsed by the system) using ONLY unauthenticated GitHub API requests.
- **FR-005**: System MUST handle GitHub API rate limiting gracefully, providing clear feedback to the user when limits are reached and preventing further requests until the limit resets.
- **FR-006**: System MUST provide an on-demand "Refresh" mechanism for users to update repository issue data, rather than performing automatic background polling.
- **FR-007**: System MUST provide a automated, single-command setup for local infrastructure (database and services).
- **FR-008**: System MUST centralize all environment-specific configurations (service URLs, access keys) in a single configuration file.
- **FR-009**: System MUST provide clear instructions and configuration for using the local Supabase CLI to allow development without external internet dependencies.
- **FR-010**: System MUST allow users to securely terminate their session (Logout).
- **FR-011**: System MUST link to existing repository data if a repository is already tracked by another user, rather than creating a duplicate global record.

### Key Entities *(include if feature involves data)*

- **User**: Represents an authenticated person. Identified by a unique `user_id` provided by Supabase. Primary attribute is the user's email address.
- **Repository**: A public GitHub repository. Attributes include the GitHub full name (e.g., "owner/repo"). Multiple users can track the same repository (Many-to-Many relationship).
- **User-Repository Link**: Connects a `User` to a `Repository`, representing the user's personal inventory.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new developer can achieve a "First Run" of the application (backend, frontend, database) in under 15 minutes.
- **SC-002**: Users can complete the login-to-dashboard flow in under 60 seconds (excluding email delivery time).
- **SC-003**: 100% of legacy GitHub OAuth code, routes, and database columns are removed from the codebase.
- **SC-004**: Data isolation is 100% verified; no user can access repositories added by another user via the API.

## Assumptions

- **Supabase Dependency**: The project assumes an active Supabase project for production authentication.
- **Public API Limits**: It is assumed that GitHub's unauthenticated API rate limits are sufficient for the expected usage, or that a single server-wide token can be configured.
- **Docker Usage**: It is assumed that all developers have Docker installed for the local environment.
- **Email Access**: Users are assumed to have immediate access to their email for magic link delivery.
