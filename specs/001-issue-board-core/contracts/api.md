# API Contracts: Issue Board Core

## Endpoints

### 1. Add Repository
- **URL**: `POST /api/repos`
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
    "owner": "owner",
    "name": "repo",
    "url": "https://github.com/owner/repo"
  }
  ```

### 2. List Repositories
- **URL**: `GET /api/repos`
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

### 3. Get Repository Issues
- **URL**: `GET /api/repos/{id}/issues`
- **Response (200 OK)**:
  ```json
  {
    "repository": "owner/repo",
    "categories": [
      {
        "label": {
          "name": "bug",
          "color": "ff0000"
        },
        "issues": [
          {
            "id": "uuid",
            "number": 123,
            "title": "Fix something",
            "url": "...",
            "state": "open"
          }
        ]
      },
      {
        "label": {
          "name": "Unlabeled",
          "color": "cccccc"
        },
        "issues": [...]
      }
    ]
  }
  ```

### 4. Refresh Issues
- **URL**: `POST /api/repos/{id}/refresh`
- **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Issues refreshed from GitHub"
  }
  ```
