# Testing Strategy

## Backend

Implemented tests cover table-driven bootstrap transitions, deterministic role assignment, secrecy serialization, content validation, bounded actor concurrency/panic recovery, HTTP, malformed/unauthenticated WebSocket traffic, a six-client WebSocket bootstrap, and in-memory reconnect projections. PostgreSQL reversibility, snapshot-tail recovery, property suites, and protocol load tests remain planned.

Required invariants: no duplicate ballot vote or role action; no early Quest resolution; exactly one ending; no cross-player private state; monotonically increasing revisions; strictly ordered event sequences; accepted command IDs execute once; published commitment matches the epilogue seed.

## Frontend

Implemented tests cover runtime configuration, protocol parsing, and React/PixiJS mount-update-dispose/remount behavior. Playwright discovers the lobby-room-match-navigation test and uses system Chrome; execution still depends on an environment that permits the preview server to bind. Store, live reconnect, visual, and accessibility suites remain planned.

Performance and resilience targets require production-like load validation. Tests must verify correctness and secrecy, not socket count alone.
