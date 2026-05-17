# IssueBoard Backend

This is the backend service for IssueBoard, written in Go. It provides a REST API to manage repository inventories and fetch/categorize GitHub issues.

---

## 📋 Prerequisites

### Go

Install Go **v1.25.0** or later from [https://go.dev/dl/](https://go.dev/dl/).

```bash
# Verify installation
go version
```

### Supabase CLI

The Supabase CLI is required to run a local Supabase instance (PostgreSQL + Auth).

**Install via Homebrew (recommended on macOS):**

```bash
brew install supabase/tap/supabase
```

**Verify:**

```bash
supabase --version
```

> **Note:** Docker Desktop must be running before starting Supabase locally. Install it from [https://www.docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop) if you haven't already.

### Docker Desktop

Supabase local dev spins up services via Docker.

- Download and install from [https://www.docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop)
- Start Docker Desktop and wait until it shows **"Running"** in the menu bar.

### GitHub Personal Access Token

A GitHub PAT is needed so the backend can call the GitHub API to fetch issues.

1. Go to [https://github.com/settings/tokens](https://github.com/settings/tokens) → **"Fine-grained tokens"** → **"Generate new token"**.
2. Grant **read-only** access to **Public Repositories** (or specific repos you want to track).
3. Copy the generated token — you'll use it as `GITHUB_TOKEN` in the `.env` file.

---

## 🏃 Running Supabase Locally

The backend connects to a local Supabase instance. These steps start that instance.

### 1. Initialize Supabase (first time only)

Run this once from the **project root** (`IssueBoard/`):

```bash
supabase init
```

### 2. Start Supabase

```bash
supabase start
```

This pulls the required Docker images and starts the following services:

| Service | Local URL |
|---|---|
| API / Auth | `http://localhost:54321` |
| Database (PostgreSQL) | `postgresql://postgres:postgres@localhost:54322/postgres` |
| Studio (dashboard) | `http://localhost:54323` |

When the command finishes it prints a summary like:

```
API URL: http://localhost:54321
DB URL: postgresql://postgres:postgres@localhost:54322/postgres
Studio URL: http://localhost:54323
anon key: eyJhbGci...
service_role key: eyJhbGci...
```

**Keep this output handy** — you'll need the `anon key` for the `.env` file.

### 3. Apply database migrations

```bash
supabase db push
```

Or, if you're inside the `backend/` directory:

```bash
supabase db push --db-url postgres://postgres:postgres@localhost:54322/postgres
```

### 4. Stop Supabase (when done)

```bash
supabase stop
```

---

## 🐛 Local Debugging

When running locally, Supabase exposes several services you can access in the browser to inspect and debug your application:

| Service | URL | What it's for |
|---|---|---|
| **Supabase Studio** | [localhost:54323](http://localhost:54323) | Full database GUI — browse tables, run SQL, inspect auth users, view logs |
| **Mailpit** (email) | [localhost:54324](http://localhost:54324) | Catches all emails sent by Supabase Auth (sign-up confirmations, magic links, password resets) |
| **Supabase Auth API** | [localhost:54321/auth/v1](http://localhost:54321/auth/v1) | REST API for authentication — useful for testing sign-up/sign-in flows via cURL |
| **Supabase REST API** | [localhost:54321/rest/v1](http://localhost:54321/rest/v1) | Auto-generated PostgREST API — allows querying the DB directly without Go |
| **PostgreSQL** | `postgresql://postgres:postgres@localhost:54322/postgres` | Direct DB connection — use any Postgres client (e.g. psql, TablePlus, DBeaver) |

### Common debugging tips

- **Emails not arriving?** Open [Mailpit](http://localhost:54324) — all emails sent during local development land there instead of a real inbox.
- **Inspect users or rows?** Open [Studio](http://localhost:54323) → **Table Editor** or the **Authentication** tab.
- **Test a query?** Open [Studio](http://localhost:54323) → **SQL Editor** to run raw SQL against the local database.
- **Auth token issues?** Use the **Authentication → Users** section in Studio to check if the user was created, or inspect JWT claims at [jwt.io](https://jwt.io).

---


## ⚙️ Configuration

Create a `.env` file inside the `backend/` directory:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:54322/postgres
SUPABASE_URL=http://localhost:54321
SUPABASE_ANON_KEY=<anon key from supabase start output>
GITHUB_TOKEN=<your GitHub personal access token>
PORT=8080

# Log level — valid values (least → most verbose):
#   error (default) | warn | info | debug | trace
# Leave unset or "error" in production.
LOG_LEVEL=debug
```

### Variable reference

| Variable | Description | How to get it |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection string for the local Supabase DB | Printed by `supabase start` as **DB URL** |
| `SUPABASE_URL` | Base URL of the local Supabase Auth/API server | Printed by `supabase start` as **API URL** (default: `http://localhost:54321`) |
| `SUPABASE_ANON_KEY` | Public anon JWT used by the auth middleware to fetch JWKS | Printed by `supabase start` as **anon key** |
| `GITHUB_TOKEN` | GitHub PAT for calling the GitHub REST API | Created in GitHub → Settings → [Personal access tokens](https://github.com/settings/tokens) |
| `PORT` | Port the backend HTTP server listens on | Any free port; default `8080` |
| `LOG_LEVEL` | Controls log verbosity | One of `error`, `warn`, `info`, `debug`, `trace` |

---

## 🚀 Setup

### 1. Install Go dependencies

```bash
cd backend/
go mod download
```

### 2. Run the server

```bash
go run cmd/main.go
```

The server starts on the port set in `.env` (default: **8080**). You should see:

```
INFO  Server starting on port 8080
```

---

## 🧪 Testing

Run all unit and integration tests:

```bash
go test ./...
```

## 📊 Test Coverage

Generate and view the HTML coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🛠️ Sample cURL Commands

### Getting an access token

All protected endpoints require a Supabase JWT in the `Authorization: Bearer <token>` header. For local development, sign up / sign in through the Supabase Auth API to get a token:

**Sign up (first time):**

```bash
curl -X POST http://localhost:54321/auth/v1/signup \
  -H "Content-Type: application/json" \
  -H "apikey: <SUPABASE_ANON_KEY>" \
  -d '{"email": "dev@example.com", "password": "password123"}'
```

**Sign in (subsequent times):**

```bash
curl -X POST http://localhost:54321/auth/v1/token?grant_type=password \
  -H "Content-Type: application/json" \
  -H "apikey: <SUPABASE_ANON_KEY>" \
  -d '{"email": "dev@example.com", "password": "password123"}'
```

Both responses include an `access_token` field. Copy that value and use it as `YOUR_ACCESS_TOKEN` in the commands below.

---

### Health Check (public)

```bash
curl http://localhost:8080/api/health
```

---

### Repositories

**Add a repository:**

```bash
curl -X POST http://localhost:8080/api/repos \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/google/go-github"}'
```

**List repositories:**

```bash
curl http://localhost:8080/api/repos \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Remove a repository:**

```bash
curl -X DELETE http://localhost:8080/api/repos/{repo_id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### Issues

**Get categorized issues for a repository:**

```bash
curl http://localhost:8080/api/repos/{repo_id}/issues \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Refresh issues from GitHub:**

```bash
curl -X POST http://localhost:8080/api/repos/{repo_id}/refresh \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### Labels

**Get all available labels for a repository:**

```bash
curl http://localhost:8080/api/repos/{repo_id}/labels/available \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get tracked labels for a repository:**

```bash
curl http://localhost:8080/api/repos/{repo_id}/labels/tracked \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Update tracked labels:**

```bash
curl -X PUT http://localhost:8080/api/repos/{repo_id}/labels/tracked \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"labels": ["bug", "good first issue", "help wanted"]}'
```

---

## 📖 API Reference

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/health` | GET | ❌ | Basic health check |
| `/api/repos` | GET | ✅ | List all repositories in the user's inventory |
| `/api/repos` | POST | ✅ | Add a repository to the inventory via GitHub URL |
| `/api/repos/{id}` | DELETE | ✅ | Remove a repository from the inventory |
| `/api/repos/{id}/issues` | GET | ✅ | Get issues grouped by tracked labels |
| `/api/repos/{id}/refresh` | POST | ✅ | Manually trigger a GitHub issue refresh |
| `/api/repos/{id}/labels/available` | GET | ✅ | Get all labels available on the GitHub repo |
| `/api/repos/{id}/labels/tracked` | GET | ✅ | Get labels currently being tracked |
| `/api/repos/{id}/labels/tracked` | PUT | ✅ | Update the set of tracked labels |

---

## 🏗️ Architecture

| Layer | Package | Responsibility |
|---|---|---|
| **Handlers** | `internal/api` | Parse HTTP requests, encode JSON responses |
| **Services** | `internal/service` | Business logic, orchestration |
| **Repositories** | `internal/repository` | Database access via `pgx` |
| **GitHub** | `internal/github` | Wraps `google/go-github` for external API calls |
| **Middleware** | `internal/api/middleware` | Auth (JWKS/JWT), logging, CORS, panic recovery |
| **Config** | `internal/config` | Loads and validates environment variables |
| **Logger** | `internal/logger` | Configures `logrus` with level from env |
