# IssueBoard Backend

This is the backend service for IssueBoard, written in Go. It provides a REST API to manage repository inventories and fetch/categorize GitHub issues.

## 📋 Prerequisites

- **Go**: v1.25.0 or later
- **PostgreSQL**: Running locally via Docker or a Supabase instance
- **GitHub Personal Access Token**: For API calls (if running without OAuth during development)

## ⚙️ Configuration

Create a `.env` file in this directory with the following variables:

```env
DATABASE_URL=postgres://user:password@localhost:5432/issueboard
SUPABASE_URL=your_supabase_url
SUPABASE_ANON_KEY=your_supabase_anon_key
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
PORT=8080
```

## 🚀 Setup

### Backend Setup
1. Install dependencies:
   ```bash
   go mod download
   ```
2. Run the server:
   ```bash
   go run cmd/main.go
   ```
   The server will start on the port specified in `.env` (default: 8080).

### Frontend Setup (Reference)
1. Navigate to the frontend directory:
   ```bash
   cd ../frontend
   ```
2. Install dependencies and start development server:
   ```bash
   npm install
   npm run dev
   ```

## 🧪 Testing

Run all unit and integration tests:
```bash
go test ./...
```

## 📊 Test Coverage

To generate and view the test coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🛠️ Sample CURL Commands

Note: Most endpoints require a `Authorization: Bearer <token>` header from Supabase.

### Health Check (Public)
```bash
curl http://localhost:8080/api/health
```

### Add a Repository
```bash
curl -X POST http://localhost:8080/api/repos \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/google/go-github"}'
```

### List Repositories
```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/repos
```

### Get Categorized Issues
```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/repos/{repo_id}/issues
```

### Refresh Issues
```bash
curl -X POST -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/repos/{repo_id}/refresh
```

## 📖 API Documentation

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/health` | GET | No | Basic health check |
| `/api/repos` | GET | Yes | List all repositories in the user's inventory |
| `/api/repos` | POST | Yes | Add a new repository to the inventory via URL |
| `/api/repos/{id}/issues` | GET | Yes | Get issues for a repository grouped by labels |
| `/api/repos/{id}/refresh` | POST | Yes | Manually trigger a refresh of issues from GitHub |

### Data Layer
- **Handlers**: Process HTTP requests and manage JSON encoding/decoding.
- **Services**: Contain business logic and coordinate between repositories and external APIs.
- **Repositories**: Handle database interactions using `pgx`.
- **GitHub Service**: Wraps the `google/go-github` client for external API interaction.
