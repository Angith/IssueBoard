# Research: Issue Board Core

## Decision: GitHub OAuth via Supabase
- **Decision**: Use Supabase Auth for GitHub OAuth integration.
- **Rationale**: Supabase provides a managed GoTrue service that handles GitHub OAuth flows seamlessly. It simplifies session management and user storage.
- **Alternatives considered**: Direct GitHub OAuth implementation (rejected for higher maintenance overhead).

## Decision: Database Connection (Go)
- **Decision**: Use `pgx` as the PostgreSQL driver for Go.
- **Rationale**: `pgx` is performant, supports modern PostgreSQL features, and integrates well with connection pooling.
- **Alternatives considered**: `lib/pq` (maintenance mode), GORM (rejected to maintain lean, explicit repository layer).

## Decision: GitHub API Interaction
- **Decision**: Use `google/go-github` library.
- **Rationale**: Most mature and well-maintained Go client for GitHub REST API.
- **Alternatives considered**: Custom HTTP client (rejected for unnecessary complexity).

## Decision: Frontend Framework
- **Decision**: Next.js (App Router).
- **Rationale**: Provides easy routing, server-side rendering for SEO (if needed later), and a robust developer experience.
- **Alternatives considered**: Vite + React (rejected as Next.js offers more built-in features for this type of app).

## Decision: Environment-Based Database Config
- **Decision**: Use `DATABASE_URL` as the primary connection string, with fallback to individual components (`DB_HOST`, etc.).
- **Rationale**: Follows 12-factor app principles and is natively supported by both Docker PG and Supabase.
- **Alternatives considered**: Hardcoded config files (rejected per constitution).
