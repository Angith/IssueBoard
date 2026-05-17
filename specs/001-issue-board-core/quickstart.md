# Quickstart: Issue Board Core

## Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- GitHub OAuth App (ClientID and ClientSecret)
- Supabase Account (or local Supabase instance)

## Local Development Setup

### 1. Backend
1. Navigate to `/backend`.
2. Copy `.env.example` to `.env`.
3. Fill in `DATABASE_URL` (points to local Docker PG).
4. Run `go mod download`.
5. Run `go run cmd/main.go`.

### 2. Frontend
1. Navigate to `/frontend`.
2. Copy `.env.local.example` to `.env.local`.
3. Run `npm install`.
4. Run `npm run dev`.

### 3. Infrastructure
1. Navigate to `/infra`.
2. Run `docker-compose up -d`.

## Environment Variables

### Backend
- `PORT`: Server port (default 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `GITHUB_CLIENT_ID`: OAuth Client ID
- `GITHUB_CLIENT_SECRET`: OAuth Client Secret
- `SUPABASE_URL`: Supabase Project URL
- `SUPABASE_ANON_KEY`: Supabase Anon Key

### Frontend
- `NEXT_PUBLIC_API_URL`: Backend API URL
- `NEXT_PUBLIC_SUPABASE_URL`: Supabase URL
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`: Supabase Anon Key
