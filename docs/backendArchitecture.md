# Backend Architecture

## Goals and boundaries

The backend is an authoritative, deterministic, secret-safe modular monolith in Go 1.26. PostgreSQL stores durable product data and match history; Redis supports only ephemeral sessions, presence, routing, reconnect leases, and rate limiting. HTTP, WebSocket, storage, logging, telemetry, time, and randomness are adapters around domain/application ports.

Domain packages never import infrastructure. The game engine imports no HTTP, WebSocket, PostgreSQL, Redis, or logger. Application code depends on domain and ports; adapters implement ports; transports invoke application use cases. Repositories contain no domain decisions and persistence records do not automatically become domain models.

## Engine and state machine

The near-pure engine handles `(MatchState, Command) -> Result`, returning revised state, public events, recipient-private events, and canonical domain events. Commands validate identity, authorization inputs, idempotency, expected revision, phase, target, payload, deadline, and rule invariants before mutation.

Canonical states are `LOBBY`, `ROLE_DEAL`, `CHAPTER_START`, `CRISIS_REVEAL`, `NOMINATION`, `COUNCIL_VOTE`, `PARTY_BUILD`, `QUEST_COMMIT`, `QUEST_REVEAL`, `EXECUTIVE_POWER`, `TRIAL_WINDOW`, `ROUND_END`, `FINAL_RECKONING`, `EPILOGUE`, and `PAUSED_RECONNECT`. GDD Figure 3's `EPILOGUE_REMATCH` conflicts with the normative table; engine, protocol, persistence, and telemetry use only `EPILOGUE`. Rematch, script selection, and departure are outgoing product destinations.

Each transition specification records entry condition, allowed commands, timeout/default, exit condition, emitted events, possible next states, and recovery. The complete allowed-command table remains normative in GDD §11.3; the vertical slice implements only lobby start.

## Room actor and command lifecycle

Exactly one goroutine owns each active match and processes a bounded mailbox sequentially. Registry synchronization protects room lookup only; no global match mutation lock exists. A transport-authenticated command enters the owning mailbox, the engine executes, and the in-memory state advances before explicit projections are delivered. Capacity errors create backpressure; slow clients have bounded outbound queues and are disconnected. Transactional event persistence is planned, not implemented in this slice.

Every command carries command, match, and player IDs, expected revision, type, payload, timestamp, and monotonic client sequence. Accepted results receive server timestamp, revision, and event sequence. Duplicate accepted IDs return the original disposition without execution.

## Persistence and recovery

The migration foundation defines PostgreSQL tables for users, profiles, rooms, matches, participants, append-only events, compressed snapshots, role statistics, reports, and asset manifests. Event mutation is rejected by a database trigger. Runtime persistence is deferred. The required future policy is a compressed snapshot every 20 events and every phase boundary; this is not yet executed by the in-memory slice.

Implemented reconnect replaces the same player's active in-memory WebSocket and returns public state plus only that player's private projection without mutating the match. Durable node-failure recovery—snapshot validation, event-tail replay, timer restoration, and owner-node routing—is deferred.

## Randomness, content, and ending integrity

A crypto-secure 32-byte seed is generated at match start and only its SHA-256 commitment is published. HMAC-derived, domain-labeled streams independently drive roles, event order, targets, and Sigil shuffle. Fixed test seeds support simulations. The seed is encrypted at rest and revealed only in `EPILOGUE` for replay verification.

Static roles, scenarios, and events are versioned JSON loaded in deterministic path/ID order. Startup rejects unknown schema versions, empty/duplicate IDs, invalid values, and broken cross-references. Stable localization keys avoid embedding final prose in rules.

## API, security, operations, and scaling

OpenAPI defines request/response endpoints; AsyncAPI defines WebSocket control envelopes. The delivery adapter derives `PublicMatchView`, the authenticated `PrivatePlayerView`, or `SpectatorView`; it never serializes `InternalMatchState`. Authentication precedes room authorization. Active secrets are encrypted and logs redact all hidden data.

Shutdown stops HTTP admission, closes registered actors, shuts down listeners within deadline, and terminates. Persistence draining is deferred with runtime persistence. Future horizontal scaling uses a match directory or consistent hashing to route a room to one node; migration and distributed actor frameworks are intentionally deferred.

Instrument command latency/rejections, mailbox depth, phases, events, snapshots, connections, reconnects, database/Redis calls, and panics with structured logs, metrics, and OpenTelemetry. Targets requiring load tests: acknowledgement p95 ≤150 ms/p99 ≤300 ms in-region, ≤1 MB per room excluding archive, reconnect ≤2 seconds, typical public events ≤4 KB, compressed snapshots ≤64 KB. Voice media is WebRTC-only and never enters this server.

## Testing

Use domain/transition unit tests, invariant/property tests, simulations, adapter integration tests, contract tests, recovery tests, and protocol load tests. Verify idempotency, strict ordering, monotonic revision, single endings, timer correctness, and zero cross-player/cross-room secret leakage.
