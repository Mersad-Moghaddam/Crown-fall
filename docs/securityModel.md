# Security Model

The implemented WebSocket boundary requires a deliberately non-production `Bearer test-<playerId>` token and binds every command's player and match IDs to the authenticated URL scope. Production identity, secure cookies/short-lived tokens, rate limiting, and reconnect leases are deferred. Engine commands validate membership, host authority, phase, revision, monotonic client sequence, and command ID; later gameplay must also validate seat and target sets.

`InternalMatchState` is never serialized. The server derives `PublicMatchView`, `PrivatePlayerView`, or restricted `SpectatorView`, then removes all other recipients from private event partitions before delivery. Active role mappings and unrevealed seeds are encrypted at rest with keys outside logs and database records.

WebSocket input is limited to 16 KiB and envelope JSON/version/identity are validated. Accepted command IDs are idempotent and conflicting reuse is rejected. Explicit localhost origin checks and bounded 32-message outbound queues are implemented; a full queue cancels the slow connection. Rate limiting, heartbeat expiry, durable presence leases, and chat moderation remain planned.

Logs, traces, browser diagnostics, and crash reports redact roles, objectives, hidden cards, seeds, tokens, and raw private payloads. Database roles follow least privilege; secrets come from runtime configuration, never source. Dependency and container scanning run in CI.

The seed commitment is published before randomness is consumed. Role selection, event ordering, target selection, and Sigil shuffling use domain-separated streams. The seed is revealed only in `EPILOGUE` so auditors can verify the commitment and replay. Spectators receive only delayed, explicitly safe projections.
