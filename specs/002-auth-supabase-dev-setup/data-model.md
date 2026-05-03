# Data Model: Auth & User Isolation

## Entities

### User
Represents an authenticated user from Supabase.
- `id` (UUID, PK): The `sub` claim from Supabase JWT.
- `email` (Text): User's email address.
- `created_at` (Timestamp): Record creation time.

### Repository
Represents a public GitHub repository. Shared across users.
- `id` (UUID, PK): Internal unique identifier.
- `github_repo_id` (BigInt, Unique): GitHub's official repository ID.
- `full_name` (Text): e.g., "google/go-github".
- `owner` (Text): e.g., "google".
- `name` (Text): e.g., "go-github".
- `url` (Text): Public URL.
- `last_fetched_at` (Timestamp): Last time issues were synced.

### User_Repository (Join Table)
Connects users to the repositories they track.
- `user_id` (UUID, FK -> User.id): The user tracking the repo.
- `repository_id` (UUID, FK -> Repository.id): The repo being tracked.
- `added_at` (Timestamp): When the user added it to their list.
- **Primary Key**: `(user_id, repository_id)`

## Relationships
- **User** ↔ **Repository**: Many-to-Many via `User_Repository`.
- **Repository** ↔ **Issue**: One-to-Many.
- **Issue** ↔ **Label**: Many-to-Many.

## Validation Rules
- `User.id` must be a valid UUID.
- `Repository.github_repo_id` must be positive.
- A user cannot add the same repository twice.
- Only repositories with a valid `full_name` and `owner` are stored.
