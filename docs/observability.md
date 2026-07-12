# Observability

Implemented observability is currently limited to structured JSON HTTP request/startup/shutdown logs with method, path, and duration. These logs intentionally exclude request bodies and private projections. Request/connection/match/command correlation and OpenTelemetry export remain planned.

Planned low-cardinality metrics cover active rooms, actors and WebSockets; command acceptance/rejection; revision conflicts; duplicates; mailbox saturation; reconnects; transitions; actor panic recovery; persistence and snapshot duration; dependency latency; and client frame-time buckets. Raw player IDs must never be metric labels.

Recover panics at process/connection boundaries, terminate the affected room safely, and alert on panic rate, persistent write lag, command latency budget breaches, reconnect spikes, queue saturation, cross-room delivery invariant failures, and readiness loss. Never log private roles, objectives, cards, or unrevealed seeds; sensitive diagnostics are explicit, audited, and disabled by default.
