# Feature Specification: Issue Board Core

**Feature Branch**: `001-issue-board-core`  
**Created**: 2026-04-26  
**Status**: Draft  
**Input**: User description: "Create a detailed feature specification for a web application called issueBoard. Feature Goal: Build a system where users can maintain an inventory of GitHub repositories and view their open issues organized dynamically by labels..."

## Overview
IssueBoard is a web application designed to help developers manage and view GitHub issues across multiple repositories in a single, unified view. Users can curate an inventory of repositories and view issues grouped dynamically by their labels, providing a clear overview of work items categorized by their type, priority, or component.

## Clarifications

### Session 2026-04-26
- Q: What level of interaction should users have with issues/repositories in the app? → A: Read-only with external links to GitHub in new tabs.
- Q: How should issues with multiple labels be handled in the board view? → A: Duplicated: The issue appears in every label section it is associated with.
- Q: Should the app support private repositories or only public ones? → A: Public & Private: Support any repository the authenticated user has access to.
- Q: Should the system implement background synchronization for issues? → A: No, remove background sync from the specification for now.
- Q: Should repository inventories be shared between users or private to individuals? → A: Individual Users: Every user has their own private inventory of repositories.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - GitHub Authentication (Priority: P1)
As a developer, I want to log in using my GitHub account so that the system can securely access my repositories and personal settings.

**Why this priority**: Authentication is the foundational entry point for the application. Without it, the system cannot identify the user or fetch their repository data securely.

**Independent Test**: Can be fully tested by attempting to log in via GitHub; delivers a secure user session and access token.

**Acceptance Scenarios**:
1. **Given** a user is on the login page, **When** they click "Login with GitHub", **Then** they are redirected to GitHub for authorization.
2. **Given** a user has authorized the app on GitHub, **When** they are redirected back, **Then** a secure session is established and they are taken to the repository inventory page.

---

### User Story 2 - Repository Inventory Management (Priority: P1)
As a user, I want to add GitHub repositories to my inventory so that I can track their issues in one place.

**Why this priority**: Curating an inventory is central to the application's goal of cross-repository issue management.

**Independent Test**: Can be tested by adding a valid repository URL and verifying it appears in the user's inventory list.

**Acceptance Scenarios**:
1. **Given** a user is on the inventory page, **When** they enter a valid GitHub repository URL and click "Add", **Then** the system validates the repository and adds it to their inventory.
2. **Given** a repository already exists in the inventory, **When** the user tries to add it again, **Then** the system prevents the duplicate entry and shows an appropriate message.
3. **Given** a user has access to a private repository, **When** they enter the private repository URL, **Then** the system validates access using their OAuth token and adds it to the inventory.
4. **Given** User A and User B both add the same repository, **When** User A views their inventory, **Then** they only see their own instance and settings for that repository.

---

### User Story 3 - Categorized Issue Board View (Priority: P1)
As a user, I want to view issues for a repository grouped by labels so that I can easily see work items organized by category.

**Why this priority**: This is the core value proposition of the application—visualizing issues dynamically based on their labels.

**Independent Test**: Can be tested by selecting a repository and verifying that issues are displayed and correctly grouped by their associated labels.

**Acceptance Scenarios**:
1. **Given** a repository with labeled issues, **When** the user views the repository board, **Then** issues are displayed in sections corresponding to their GitHub labels and link to GitHub.
2. **Given** an issue has multiple labels, **When** the user views the board, **Then** the issue appears in every label section it is associated with.
3. **Given** an issue has no labels, **When** the user views the board, **Then** the issue is grouped under a default "Unlabeled" category.

---

### User Story 4 - Issue Caching & On-Demand Refresh (Priority: P2)
As a user, I want the issue board to load quickly using cached data and have an option to refresh manually so that I have a performant and accurate view of my work.

**Why this priority**: Performance is critical for a good user experience, and on-demand refreshing allows users to get the latest data when needed without constant background overhead.

**Independent Test**: Can be tested by loading a board, verifying it uses cached data first, and then triggering a manual refresh.

**Acceptance Scenarios**:
1. **Given** issues are already stored in the database, **When** a user opens the board, **Then** the system displays the cached issues immediately.
2. **Given** the user triggers a refresh on the board, **When** the request is made, **Then** the system fetches the latest issues from GitHub, updates the local cache, and refreshes the view.

---

## Edge Cases
- **Invalid GitHub URL**: User enters a string that is not a valid GitHub repository URL.
- **Private Repo Access**: User adds a private repository they don't have permission to access.
- **Empty Repository**: A repository is added that has no issues.
- **GitHub Rate Limiting**: The system reaches GitHub's API rate limit during an issue fetch.
- **Database Connection Failure**: The application loses connection to the PostgreSQL database.
- **Issue with many labels**: An issue will appear in all relevant label sections as per duplication strategy.

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST authenticate users via GitHub OAuth.
- **FR-002**: System MUST allow users to add repositories via GitHub URL and store metadata.
- **FR-003**: System MUST prevent duplicate repositories in a user's inventory.
- **FR-004**: System MUST fetch and store issues and labels from the GitHub REST API.
- **FR-005**: System MUST group issues dynamically by their labels in the UI.
- **FR-006**: System MUST provide an "Unlabeled" category for issues without labels.
- **FR-007**: System MUST implement a caching layer (PostgreSQL) to minimize external API calls.
- **FR-008**: System MUST support environment-based database configuration (local Docker vs. Supabase).
- **FR-009**: Backend MUST follow a layered architecture (Handlers -> Services -> Repositories).
- **FR-010**: System MUST provide direct links to the original GitHub repository and issue (opening in a new tab).
- **FR-011**: System interaction for issues and repositories is READ-ONLY.
- **FR-012**: System MUST duplicate issues across multiple label sections if multiple labels are present.
- **FR-013**: System MUST support both public and private repositories accessible by the authenticated user.
- **FR-014**: System MUST provide a mechanism to manually refresh issues for a repository on-demand.
- **FR-015**: System MUST ensure repository inventories are private to each individual user.

### Key Entities
- **User**: Represents a registered user. Attributes: ID, GitHub ID, OAuth Token, Username.
- **Repository**: A GitHub repository tracked by a user. Attributes: ID, GitHub Repo ID, Full Name, Owner, URL, UserID.
- **Issue**: A GitHub issue. Attributes: ID, GitHub Issue ID, Title, Body, Number, State, RepositoryID.
- **Label**: A GitHub label. Attributes: ID, Name, Color, RepositoryID.
- **Relationships**: User -> has many -> Repositories; Repository -> has many -> Issues; Issue -> has many -> Labels.

## Success Criteria *(mandatory)*

### Measurable Outcomes
- **SC-001**: Users can log in and view their repository inventory in under 2 seconds.
- **SC-002**: The issue board loads cached issues in under 500ms.
- **SC-003**: 100% of repositories added are validated for existence and permissions before being saved.
- **SC-004**: System stays within GitHub API rate limits for standard usage via effective caching.
- **SC-005**: All database migrations run successfully and consistently in both local and production environments.

## Assumptions
- **Connectivity**: Users have a stable internet connection for GitHub API interactions.
- **GitHub API**: The GitHub REST API is available and follows standard response formats.
- **Browser Support**: The frontend targets modern evergreen browsers (Chrome, Firefox, Safari, Edge).
- **OAuth Permissions**: The app will request standard `repo` or `public_repo` scopes as needed.
- **Data Persistence**: Issues and labels are stored in PostgreSQL; full history of issue comments is out of scope for v1.
