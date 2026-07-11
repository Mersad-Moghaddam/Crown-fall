# Crownfall Backend

The Go backend is the authoritative owner of match state, hidden information, validation, deterministic randomness, timers, projections, history, and outcomes. Clients submit intentions only.

## Architecture

This is a hexagonal modular monolith. Domain and engine packages have no infrastructure imports. Application code invokes ports; adapters implement HTTP, WebSocket, persistence, and operational integration. One bounded-mailbox goroutine serializes commands for each active match. PostgreSQL is durable truth; Redis is reserved for sessions, presence, routing, reconnect leases, and rate limiting.

The engine accepts a state and command and returns revised state plus public, recipient-private, and domain events. It must run without a database, Redis, network, logger, or real clock. Randomness publishes a SHA-256 seed commitment and derives separate HMAC streams for roles, event order, targets, and Sigils; reveal occurs only in `EPILOGUE`.

## Environment and commands

Copy `.env.example` to `.env` and configure HTTP, database, Redis, log level, and shutdown timeout.

```sh
go run ./cmd/server
go run ./cmd/simulate
go run ./cmd/migrate up
go test ./...
go vet ./...
go build ./cmd/server
```

The migration command is currently an explicit adapter boundary; SQL files are applied by the deployment migration tool. WebSocket contracts live in `api/asyncapi.yaml`.

## Extending the system

- Command: define the typed payload in the engine, validate identity/revision/phase, return explicit events, update AsyncAPI, and add table-driven tests.
- Transition: use the exact GDD identifier and state table, define entry/exit/timeout/recovery behavior, and add invariant tests.
- Role or scenario: add versioned JSON with stable IDs and localization keys, extend validation, and add fixtures. Never put rules in transport or repository code.
- Repository: define a port near the application use case and implement it in an adapter; do not reuse persistence records as domain models by default.

If startup fails, check `CROWNFALL_HTTP_ADDRESS`. Database and Redis are not required by the current in-memory slice. A WebSocket origin must match the configured development policy.
