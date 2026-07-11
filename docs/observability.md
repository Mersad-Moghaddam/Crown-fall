# Observability

Emit structured JSON logs and OpenTelemetry-compatible traces/metrics across HTTP upgrade, authenticated connection, room mailbox, command validation, transition, projection, persistence, and snapshot boundaries. Correlate with request, connection, match, command, and trace IDs; use pseudonymous player-safe identifiers.

Measure active rooms/connections, mailbox depth, command latency and rejection reason, phase duration, reconnect rate/restore latency, snapshot duration, event-log lag, database/Redis latency, client frame-time buckets, role win rates, and ending distribution. Categorize validation, authorization, conflict, dependency, capacity, and internal failures.

Recover panics at process/connection boundaries, terminate the affected room safely, and alert on panic rate, persistent write lag, command latency budget breaches, reconnect spikes, queue saturation, cross-room delivery invariant failures, and readiness loss. Never log private roles, objectives, cards, or unrevealed seeds; sensitive diagnostics are explicit, audited, and disabled by default.
