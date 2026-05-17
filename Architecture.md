# 🧭 Architecture Overview

This document describes the system architecture and data flow for IssueBoard.

The application allows users to:
- Authenticate via Supabase Magic Link (passwordless, email-based)
- Add GitHub repositories to their personal inventory
- Fetch available labels from those repositories and select which ones to track
- View issues grouped by their tracked labels
- Manually trigger a refresh to pull the latest issues from GitHub

---

# 🏗️ High-Level Architecture

```mermaid
flowchart TD

U[User Browser]
FE["Frontend — React / Next.js"]
AUTH["Auth — Supabase Magic Link"]
BE["Backend API — Go"]
DB["PostgreSQL — Supabase"]
GH[GitHub API]

U -->|Login / UI Actions| FE

FE -->|Magic Link request| AUTH
AUTH -->|Send magic link email| U
U -->|Click email link| AUTH
AUTH -->|JWT access_token| FE

FE -->|API Requests + JWT| BE

BE -->|Verify JWT via JWKS| AUTH
BE -->|Read / Write| DB
BE -->|Fetch repos / issues / labels| GH

BE -->|JSON Response| FE
FE -->|Render UI| U
```

> **Note:** There is no background worker or in-memory cache. Issues are fetched from GitHub on demand when a refresh is triggered, and stored in PostgreSQL for subsequent reads.

---

# 🔐 Authentication Flow (Supabase Magic Link)

Authentication is fully passwordless. Users enter their email, receive a one-time magic link, and click it to get a JWT access token issued by Supabase.

```mermaid
sequenceDiagram
    participant U as User
    participant FE as Frontend
    participant AUTH as Supabase Auth
    participant DB as Database

    U->>FE: Enter email address
    FE->>AUTH: POST /auth/v1/otp (email)
    AUTH-->>U: Send magic link email
    U->>AUTH: Click magic link
    AUTH-->>FE: Return JWT access_token
    FE->>BE: All subsequent requests with Bearer token
    BE->>AUTH: Validate JWT via JWKS endpoint
    BE->>DB: Upsert user record (on first request)
```

> In local development, magic link emails are intercepted by **Mailpit** at [localhost:54324](http://localhost:54324) — no real email is sent.

---

# 📦 Add Repository Flow

Users add a repository by providing its GitHub URL. The backend validates it exists on GitHub and stores the metadata.

```mermaid
sequenceDiagram
    participant U as User
    participant FE as Frontend
    participant BE as Backend
    participant GH as GitHub API
    participant DB as Database

    U->>FE: Enter GitHub repo URL
    FE->>BE: POST /api/repos {url}
    BE->>GH: Fetch repo metadata (owner, name, description)
    GH-->>BE: Repo metadata
    BE->>DB: Store repository record
    BE-->>FE: Created repository response
    FE-->>U: Repo appears in inventory
```

---

# 🏷️ Label Selection Flow

Before issues can be displayed, the user must select which labels to track. This is a one-time setup per repository that can be updated any time.

```mermaid
sequenceDiagram
    participant U as User
    participant FE as Frontend
    participant BE as Backend
    participant GH as GitHub API
    participant DB as Database

    U->>FE: Open label configuration for a repo
    FE->>BE: GET /api/repos/{id}/labels/available
    BE->>GH: Fetch all labels for repo
    GH-->>BE: Labels list
    BE-->>FE: Available labels
    FE-->>U: Show label picker

    U->>FE: Select labels to track
    FE->>BE: PUT /api/repos/{id}/labels/tracked {labels}
    BE->>DB: Save tracked labels for repo
    BE-->>FE: Success
    FE-->>U: Labels saved
```

> The user can also view currently tracked labels at any time via `GET /api/repos/{id}/labels/tracked`.

---

# 🐞 Fetch Issues Flow

Issues are fetched from the database and grouped by the user's tracked labels. If no cached issues exist for the tracked labels, a GitHub refresh is triggered automatically.

```mermaid
sequenceDiagram
    participant U as User
    participant FE as Frontend
    participant BE as Backend
    participant DB as Database
    participant GH as GitHub API

    U->>FE: View issues for a repository
    FE->>BE: GET /api/repos/{id}/issues
    BE->>DB: Load tracked labels for repo

    alt No tracked labels configured
        BE-->>FE: Empty board (is_tracking_configured: false)
        FE-->>U: Prompt to configure labels first
    else Tracked labels exist
        BE->>DB: Load cached issues for repo
        BE->>BE: Filter issues by tracked labels

        alt No cached issues for tracked labels
            BE->>GH: Fetch issues filtered by tracked labels
            GH-->>BE: Issues + label metadata
            BE->>DB: Upsert issues and sync labels
        end

        BE->>BE: Group filtered issues by label
        BE-->>FE: IssueBoard {repository, categories[]}
        FE-->>U: Render issue board grouped by label
    end
```

---

# 🔄 Manual Refresh Flow

Users can trigger a refresh at any time to pull the latest issues from GitHub for a repository.

```mermaid
sequenceDiagram
    participant U as User
    participant FE as Frontend
    participant BE as Backend
    participant DB as Database
    participant GH as GitHub API

    U->>FE: Click "Refresh"
    FE->>BE: POST /api/repos/{id}/refresh
    BE->>DB: Load tracked labels for repo

    alt No tracked labels configured
        BE-->>FE: Success (no-op)
    else Tracked labels exist
        BE->>GH: Fetch issues filtered by tracked labels
        GH-->>BE: Issues + label metadata
        BE->>DB: Upsert issues, sync labels
        BE-->>FE: Success
        FE-->>U: Issues refreshed
    end
```

---
