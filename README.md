# IssueBoard

IssueBoard is a web application designed to help developers manage and view GitHub issues across multiple repositories in a single, unified view. Users can curate an inventory of repositories and view issues grouped dynamically by their labels, providing a clear overview of work items categorized by their type, priority, or component.

## 🏗️ Architecture Summary

The project follows a decoupled client-server architecture:

- **Frontend**: Next.js 16 (React 19) application styled with Tailwind CSS, using Supabase for authentication.
- **Backend**: Go 1.25 REST API following a layered architecture (Handlers -> Services -> Repositories).
- **Database**: PostgreSQL (Supabase or local Docker) for caching GitHub metadata, issues, and labels.
- **Authentication**: GitHub OAuth via Supabase Auth.

For more details, see [architecture.md](./architecture.md).

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: v1.25.0 or later
- **Node.js**: v20.x or later (includes npm)
- **Docker & Docker Compose**: For running the local database
- **GitHub Account**: To create a GitHub OAuth App for authentication

## ⚙️ Configurations

### Backend Configuration
Create a `.env` file in the `backend/` directory:
```env
DATABASE_URL=postgres://user:password@localhost:5432/issueboard
SUPABASE_URL=your_supabase_url
SUPABASE_ANON_KEY=your_supabase_anon_key
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
PORT=8080
```

### Frontend Configuration
Create a `.env.local` file in the `frontend/` directory:
```env
NEXT_PUBLIC_SUPABASE_URL=your_supabase_url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 🚀 Setup

### 1. Database Setup
Start the local PostgreSQL database using Docker Compose:
```bash
docker-compose -f infra/docker-compose.yml up -d
```

### 2. Backend Setup
```bash
cd backend
go mod download
go run cmd/main.go
```
The backend will be available at `http://localhost:8080`.

### 3. Frontend Setup
```bash
cd frontend
npm install
npm run dev
```
The frontend will be available at `http://localhost:3000`.

## 🧪 Testing

### Backend
Run all backend tests:
```bash
cd backend
go test ./...
```

### Frontend
Run frontend linting:
```bash
cd frontend
npm run lint
```
*(Note: Automated tests for frontend are currently being implemented.)*

## 📊 Test Coverage

### Backend Coverage
To view backend test coverage:
```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---
For more detailed information, please refer to the specific READMEs:
- [Backend README](./backend/README.md)
- [Frontend README](./frontend/README.md)
