<!--
Sync Impact Report
- Version change: 1.0.0 → 1.1.0
- List of modified principles:
  - V. Security & Authentication: Updated to mandate Supabase Magic Link and JWT validation.
- Added sections: N/A
- Removed sections: N/A
- Templates requiring updates:
  - .specify/templates/plan-template.md (✅ aligns)
  - .specify/templates/spec-template.md (✅ aligns)
  - .specify/templates/tasks-template.md (✅ aligns)
- Follow-up TODOs: N/A
-->

# IssueBoard Constitution

## Core Principles

### I. Architecture & Separation
The system must maintain a strict separation between frontend and backend. The backend must follow a modular, layered architecture (handlers, services, repositories) and remain stateless where possible.

### II. API-First Design
All features must start with API contract definitions before implementation. Communication between frontend and backend must use RESTful conventions and JSON schemas.

### III. Clean Code & Idiomatic Standards
Code must prioritize readability, small functions, and meaningful naming. Backend code must follow idiomatic Go practices. Frontend must use consistent React patterns (functional components, hooks). Comments are required only for non-trivial logic.

### IV. Testing & Quality Gates
The backend must have a minimum of 70% unit test coverage for business logic, with external dependencies mocked. The frontend must enforce component-level validation and type safety.

### V. Security & Authentication
Authentication must be handled via Supabase Magic Link (email-based login). Backend requests must be secured by validating Supabase JWTs on every request. Secrets must be managed via environment variables, and all external integrations must be abstracted via service layers.

## Linting & Tooling
- **Backend**: Must pass `gofmt` and `golangci-lint` checks.
- **Frontend**: Must pass `ESLint` and `Prettier` formatting and linting.
- All code must pass these checks in CI before being merged.

## Documentation & Workflow
- **Development**: All features must be spec-driven. No direct commits to the `main` branch are allowed.
- **Reviews**: Pull Request reviews are mandatory for all changes.
- **Docs**: Each module must have basic documentation, and APIs must be documented using OpenAPI/Swagger.

## Governance
This constitution supersedes all other development practices in this project. Amendments require documentation in a PR, approval from project leads, and a version bump. All PR reviews must verify compliance with these principles.

**Version**: 1.1.0 | **Ratified**: 2026-04-26 | **Last Amended**: 2026-05-02
