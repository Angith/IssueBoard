# API Contracts: Auth & User-Scoped Data

All endpoints require an `Authorization: Bearer <Supabase_JWT>` header.

## Endpoints

### 1. Add Repository
- **URL**: `POST /api/repos`
- **Header**: `Authorization: Bearer <JWT>`
- **Request Body**:
  ```json
  {
    "url": "https://github.com/owner/repo"
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "id": "uuid",
    "full_name": "owner/repo",
    "url": "https://github.com/owner/repo"
  }
  ```
- **Note**: If the repository already exists globally, it is linked to the current user.

### 2. List Repositories
- **URL**: `GET /api/repos`
- **Header**: `Authorization: Bearer <JWT>`
- **Response (200 OK)**:
  ```json
  [
    {
      "id": "uuid",
      "full_name": "owner/repo",
      "owner": "owner",
      "name": "repo"
    }
  ]
  ```
- **Note**: Only returns repositories linked to the authenticated user.

### 3. Get Repository Issues
- **URL**: `GET /api/repos/{id}/issues`
- **Header**: `Authorization: Bearer <JWT>`
- **Response (200 OK)**: (Same structure as original spec)

### 4. Refresh Issues (On-Demand)
- **URL**: `POST /api/repos/{id}/refresh`
- **Header**: `Authorization: Bearer <JWT>`
- **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Issues updated from GitHub"
  }
  ```
- **Error (429 Too Many Requests)**:
  ```json
  {
    "error": "GitHub API rate limit reached. Please try again later."
  }
  ```

### 5. Authentication (Client-Side only)
Authentication is handled via the Supabase client library:
- `supabase.auth.signInWithOtp({ email })`
- `supabase.auth.onAuthStateChange()`
