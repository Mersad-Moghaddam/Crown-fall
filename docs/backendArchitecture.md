# Backend Architecture

## Goals and boundaries

The backend is an authoritative, deterministic, secret-safe modular monolith in Go 1.26. PostgreSQL stores durable product data and match history; Redis supports only ephemeral sessions, presence, routing, reconnect leases, and rate limiting. HTTP, WebSocket, storage, logging, telemetry, time, and randomness are adapters around domain/application ports.

Domain packages never import infrastructure. The game engine imports no HTTP, WebSocket, PostgreSQL, Redis, or logger. Application code depends on domain and ports; adapters implement ports; transports invoke application use cases. Repositories contain no domain decisions and persistence records do not automatically become domain models.

## Engine and state machine

The near-pure engine handles `(MatchState, Command) -> Result`, returning revised state, public events, recipient-private events, and canonical domain events. Commands validate identity, authorization inputs, idempotency, expected revision, phase, target, payload, deadline, and rule invariants before mutation.

Canonical states are `LOBBY`, `ROLE_DEAL`, `CHAPTER_START`, `CRISIS_REVEAL`, `NOMINATION`, `COUNCIL_VOTE`, `PARTY_BUILD`, `QUEST_COMMIT`, `QUEST_REVEAL`, `EXECUTIVE_POWER`, `TRIAL_WINDOW`, `ROUND_END`, `FINAL_RECKONING`, `EPILOGUE`, and `PAUSED_RECONNECT`. The GDD Figure 3 discrepancy is resolved in favor of normative `EPILOGUE` and remains documented in the glossary.

Each transition specification records entry condition, allowed commands, timeout/default, exit condition, emitted events, possible next states, and recovery. The complete allowed-command table remains normative in GDD §11.3; the vertical slice implements only lobby start.

## Room actor and command lifecycle

Exactly one goroutine owns each active match and processes a bounded mailbox sequentially. Registry synchronization protects room lookup only; no global match mutation lock exists. A transport-authenticated command enters the owning mailbox, the application loads idempotency/recovery context, the engine executes, durable events and required snapshot commit transactionally, the in-memory state advances, and projections are delivered. Capacity errors create backpressure; slow clients have bounded outbound queues and are disconnected/resynchronized.

Every command carries command, match, and player IDs, expected revision, type, payload, timestamp, and monotonic client sequence. Accepted results receive server timestamp, revision, and event sequence. Duplicate accepted IDs return the original disposition without execution.

## Persistence and recovery

PostgreSQL stores users, profiles, rooms, matches, participants, results, append-only events, compressed snapshots, replay/audit metadata, idempotency, role statistics, conduct reports, and asset manifests. Events have an immutable unique sequence per match. A compressed snapshot is stored every 20 events and every phase boundary; this is targeted event history, not product-wide event sourcing.

On node failure, routing directs reconnects to the owner or recovery node. Recovery loads the latest valid snapshot, verifies its sequence/content version, replays the canonical tail through deterministic reducers, restores timers using injected clock policy, and resumes the actor. Reconnect returns a recipient projection plus missing tail within the ≤2 second target.

## Randomness, content, and ending integrity

A crypto-secure 32-byte seed is generated at match start and only its SHA-256 commitment is published. HMAC-derived, domain-labeled streams independently drive roles, event order, targets, and Sigil shuffle. Fixed test seeds support simulations. The seed is encrypted at rest and revealed only in `EPILOGUE` for replay verification.

Static roles, scenarios, and events are versioned JSON loaded in deterministic path/ID order. Startup rejects unknown schema versions, empty/duplicate IDs, invalid values, and broken cross-references. Stable localization keys avoid embedding final prose in rules.

## API, security, operations, and scaling

OpenAPI defines request/response endpoints; AsyncAPI defines WebSocket control envelopes. The delivery adapter derives `PublicMatchView`, the authenticated `PrivatePlayerView`, or `SpectatorView`; it never serializes `InternalMatchState`. Authentication precedes room authorization. Active secrets are encrypted and logs redact all hidden data.

Shutdown stops admission, drains bounded mailboxes within deadline, persists accepted events/snapshots, closes sockets, and terminates. Future horizontal scaling uses a match directory or consistent hashing to route a room to one node; migration and distributed actor frameworks are intentionally deferred.

Instrument command latency/rejections, mailbox depth, phases, events, snapshots, connections, reconnects, database/Redis calls, and panics with structured logs, metrics, and OpenTelemetry. Targets requiring load tests: acknowledgement p95 ≤150 ms/p99 ≤300 ms in-region, ≤1 MB per room excluding archive, reconnect ≤2 seconds, typical public events ≤4 KB, compressed snapshots ≤64 KB. Voice media is WebRTC-only and never enters this server.

## Testing

Use domain/transition unit tests, invariant/property tests, simulations, adapter integration tests, contract tests, recovery tests, and protocol load tests. Verify idempotency, strict ordering, monotonic revision, single endings, timer correctness, and zero cross-player/cross-room secret leakage.
