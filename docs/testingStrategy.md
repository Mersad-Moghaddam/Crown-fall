# Testing Strategy

## Backend

Use table-driven domain and transition tests, property-style invariant tests, deterministic simulations, repository integration tests against PostgreSQL, WebSocket contract tests, snapshot-plus-tail recovery tests, and protocol load tests for join bursts, voting, Quest commits, and reconnect storms.

Required invariants: no duplicate ballot vote or role action; no early Quest resolution; exactly one ending; no cross-player private state; monotonically increasing revisions; strictly ordered event sequences; accepted command IDs execute once; published commitment matches the epilogue seed.

## Frontend

Use component, hook, store, and protocol parser tests; React/PixiJS lifecycle tests; visual smoke tests; Playwright journeys; reconnection and accessibility tests. Contract fixtures must contain only the projection a real recipient receives.

Performance and resilience targets require production-like load validation. Tests must verify correctness and secrecy, not socket count alone.
