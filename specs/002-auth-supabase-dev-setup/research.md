# Research: Supabase Authentication and Dev Experience

## Decision: JWT Verification in Go
- **Choice**: Use `golang-jwt/jwt` for manual validation using the Supabase JWT Secret.
- **Rationale**: While `gotrue-go` is available, manual JWT validation is lightweight, has fewer dependencies, and is standard for verifying Supabase tokens in custom backends. It allows for easy extraction of the `sub` (user_id) claim.
- **Alternatives considered**: 
    - `gotrue-go`: Rejected for being more complex than needed for simple token verification.
    - Remote validation: Rejected to avoid external network calls on every request.

## Decision: Local Development Infrastructure
- **Choice**: Official Supabase CLI for Auth and local PostgreSQL.
- **Rationale**: Provides the most accurate representation of the production environment, including the Auth emulator and local database migrations. 
- **Alternatives considered**: 
    - Generic Docker Postgres + Mock Auth: Rejected as it doesn't test the Magic Link flow properly.

## Decision: Environment Variable Management
- **Choice**: `joho/godotenv` for loading `.env` and `caarlos0/env` for parsing into a structured config.
- **Rationale**: Industry standard for Go projects. `godotenv` handles the file loading, while `caarlos0/env` provides type safety and default values.
- **Alternatives considered**: 
    - `os.Getenv` manually: Rejected as it becomes messy with many variables and lacks type safety.

## Decision: GitHub API Rate Limit Handling
- **Choice**: Exponential backoff and "429 Too Many Requests" response to the frontend.
- **Rationale**: Since we are unauthenticated, we must respect the 60 req/hr limit. The frontend will show a "Rate limit reached" state.
- **Alternatives considered**: 
    - Single server token: Kept as an optional environment variable for deployment but defaulted to unauthenticated for simplicity.
