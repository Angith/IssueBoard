# Quickstart: Local Development Setup

Follow these steps to set up the project locally with Supabase Auth and Docker.

## Prerequisites
- [Docker](https://www.docker.com/) and Docker Compose
- [Go 1.25+](https://go.dev/)
- [Node.js 22+](https://nodejs.org/)
- [Supabase CLI](https://supabase.com/docs/guides/cli)

## 1. Infrastructure Setup
Start the local PostgreSQL and Supabase services:
```bash
# In the project root
supabase start
```
This command will provide you with local environment variables (`API URL`, `anon key`, etc.).

## 2. Environment Configuration
Create `.env` files for both backend and frontend.

### Backend (`backend/.env`)
```env
DATABASE_URL=postgres://postgres:postgres@localhost:54322/postgres
SUPABASE_URL=http://localhost:54321
SUPABASE_ANON_KEY=your-local-anon-key
SUPABASE_JWT_SECRET=your-local-jwt-secret
PORT=8080
```

### Frontend (`frontend/.env.local`)
```env
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-local-anon-key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 3. Run the Application

### Backend
```bash
cd backend
go run cmd/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## 4. Local Authentication
1. Open `http://localhost:3000/login`.
2. Enter your email address.
3. Check the Supabase local dashboard (`http://localhost:54323/monitor`) to "catch" the magic link email.
4. Click the link to log in.
