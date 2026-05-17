# IssueBoard

IssueBoard is a web application that helps developers manage and triage GitHub issues across multiple repositories in a single, unified view. Users authenticate with their email, curate a personal inventory of GitHub repositories, and view issues dynamically grouped by labels — making it easy to track work across projects without switching tabs.

---

## ✨ Features

- **Multi-repo inventory** — add any public GitHub repository to your personal dashboard
- **Label-based categorisation** — choose which labels to track; issues are grouped accordingly
- **On-demand refresh** — pull the latest issues from GitHub whenever you need them
- **Passwordless auth** — sign in via Supabase Magic Link (email-based, no passwords)

---

## 🏗️ Architecture

The project follows a decoupled client-server architecture:

| Layer | Technology | Role |
|---|---|---|
| **Frontend** | Next.js 16 (React 19) | UI, routing, Supabase auth client |
| **Backend** | Go 1.25 REST API | Business logic, GitHub API integration |
| **Database** | PostgreSQL via Supabase | Persistent storage for repos, issues, labels |
| **Auth** | Supabase (Magic Link) | Passwordless email authentication |
| **External API** | GitHub REST API | Fetches repository metadata and issues |

For a deeper dive into flows and diagrams, see [Architecture.md](./Architecture.md).

---

## 📁 Project Structure

```
IssueBoard/
├── backend/        # Go REST API
├── frontend/       # Next.js application
├── architecture.md # System architecture and data flow diagrams
└── specs/          # Feature specifications
```

---

## 🚀 Getting Started

Each service has its own detailed setup guide:

- **[Backend README](./backend/README.md)** — Go server setup, Supabase local dev, environment variables, API reference
- **[Frontend README](./frontend/README.md)** — Next.js setup, environment variables, available scripts

---
