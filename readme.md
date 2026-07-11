# Crownfall

Crownfall is a browser-based fantasy social-deduction game for 6–10 players. This repository is the pre-production architecture baseline: a React/PixiJS client and an authoritative Go match server, developed and deployed independently.

## Status

The repository contains a minimal vertical slice, contracts, database foundations, and architectural documentation. It does not contain the complete role or scenario library, production authentication, matchmaking, or voice integration. The [game design document](crownfallGameDesignDocument.md) is the source of truth for game behavior.

## Repository

- `frontend/`: React application shell and PixiJS table renderer.
- `backend/`: authoritative game engine, room actors, HTTP/WebSocket transports, contracts, and migrations.
- `docs/`: architecture, security, testing, operations, decisions, and diagrams.
- `infrastructure/`: local, monitoring, and deployment extension points.

Frontend and backend do not import each other's source. Their only shared boundary is the versioned OpenAPI/AsyncAPI contract under `backend/api/`.

## Prerequisites

- Node.js 24 LTS and npm
- Go 1.26 with toolchain 1.26.5
- Docker with Compose (for PostgreSQL and Redis)
- GNU Make

## Start locally

```sh
cp frontend/.env.example frontend/.env
cp backend/.env.example backend/.env
make setup
make dev
```

Or run the four-service stack:

```sh
docker compose up --build
```

The frontend is served at `http://localhost:5173`; backend health and readiness are at `http://localhost:8080/healthz` and `/readyz`.

## Commands

| Task | Command |
|---|---|
| Tests | `make test` |
| Lint and vet | `make lint` |
| Production builds | `make build` |
| Frontend only | `make devFrontend` |
| Backend only | `make devBackend` |
| Migrations | `make migrateUp` / `make migrateDown` |
| Simulation | `make simulate` |

See [development workflow](docs/developmentWorkflow.md), [frontend architecture](docs/frontendArchitecture.md), [backend architecture](docs/backendArchitecture.md), and [security model](docs/securityModel.md).

## Contribution and security

Read `agents.md` and the relevant architecture document before editing. Game-rule changes require tests and GDD review; contract changes require schema updates. Never report suspected secret-state exposure in a public issue—use the private project security channel.

## License and roadmap

This is proprietary software; see `license`. The next milestones are durable room recovery, authentication, full scenario content, playtest simulation, and a separately reviewed voice-plane integration.
