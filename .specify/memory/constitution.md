<!--
Sync Impact Report
- Version change: Initial → 1.0.0
- List of modified principles:
  - [PRINCIPLE_1_NAME] → I. Architecture & Separation
  - [PRINCIPLE_2_NAME] → II. API-First Design
  - [PRINCIPLE_3_NAME] → III. Clean Code & Idiomatic Standards
  - [PRINCIPLE_4_NAME] → IV. Testing & Quality Gates
  - [PRINCIPLE_5_NAME] → V. Security & Authentication
- Added sections: Linting & Tooling, Documentation & Workflow
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
GitHub OAuth must be securely implemented, ensuring tokens are never exposed to the frontend. Secrets must be managed via environment variables, and all external integrations must be abstracted via service layers.

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

**Version**: 1.0.0 | **Ratified**: 2026-04-26 | **Last Amended**: 2026-04-26
