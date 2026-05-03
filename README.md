# IssueBoard

IssueBoard is a web application designed to help developers manage and view GitHub issues across multiple repositories in a single, unified view. Users can curate an inventory of repositories and view issues grouped dynamically by their labels.

## 🏗️ Architecture Summary

The project follows a decoupled client-server architecture:

- **Frontend**: Next.js 16 (React 19) using Supabase for authentication.
- **Backend**: Go 1.25 REST API following a layered architecture.
- **Database**: PostgreSQL (Supabase) for caching metadata.
- **Authentication**: Supabase Magic Link (email-based, passwordless).
- **External Integration**: Unauthenticated GitHub REST API.

For more details, see [architecture.md](./architecture.md).

## 📋 Prerequisites

- **Go**: v1.25.0+
- **Node.js**: v22+
- **Docker**: For local services
- **Supabase CLI**: For local Auth and DB development

## 🚀 Quick Start (Local Setup)

### 1. Start Infrastructure
```bash
supabase start
```
This starts local Postgres, Auth (GoTrue), and other services. It will output your local API URL and Keys.

### 2. Configure Environment
Create `.env` files based on the output of `supabase start`.

**Backend (`backend/.env`)**
```env
DATABASE_URL=postgres://postgres:postgres@localhost:54322/postgres
SUPABASE_URL=http://localhost:54321
SUPABASE_ANON_KEY=your_local_anon_key
SUPABASE_JWT_SECRET=your_local_jwt_secret
PORT=8080
```

**Frontend (`frontend/.env.local`)**
```env
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_local_anon_key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 3. Run Backend
```bash
cd backend
go run cmd/main.go
```

### 4. Run Frontend
```bash
cd frontend
npm install
npm run dev
```

## 🧪 Testing

### Backend
```bash
cd backend
go test ./...
```

### Frontend
```bash
cd frontend
npm run lint
```

---
For more detailed information, see [quickstart.md](./specs/002-auth-supabase-dev-setup/quickstart.md).
