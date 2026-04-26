# Data Model: Issue Board Core

## Entities

### User
- `id`: UUID (Primary Key)
- `github_id`: String (Unique)
- `username`: String
- `email`: String
- `oauth_token`: String (Encrypted)
- `created_at`: Timestamp
- `updated_at`: Timestamp

### Repository
- `id`: UUID (Primary Key)
- `user_id`: UUID (Foreign Key -> User.id)
- `github_repo_id`: BigInt (Unique for user)
- `full_name`: String (e.g., "owner/repo")
- `owner`: String
- `name`: String
- `url`: String
- `created_at`: Timestamp

### Issue
- `id`: UUID (Primary Key)
- `repository_id`: UUID (Foreign Key -> Repository.id)
- `github_issue_id`: BigInt (Unique for repo)
- `number`: Integer
- `title`: String
- `body`: Text
- `state`: String (open/closed)
- `url`: String
- `created_at`: Timestamp
- `updated_at`: Timestamp (from GitHub)

### Label
- `id`: UUID (Primary Key)
- `repository_id`: UUID (Foreign Key -> Repository.id)
- `name`: String
- `color`: String (hex)
- `description`: String

### IssueLabel (Join Table)
- `issue_id`: UUID (Foreign Key -> Issue.id)
- `label_id`: UUID (Foreign Key -> Label.id)

## Relationships
- A **User** has many **Repositories**.
- A **Repository** has many **Issues**.
- A **Repository** has many **Labels**.
- An **Issue** can have many **Labels**.
- A **Label** can be on many **Issues**.
