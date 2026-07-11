# Security Model

The edge authenticates HTTP sessions and WebSocket upgrades. Production sessions use secure, HTTP-only cookies or short-lived bearer material; reconnect tokens are short-lived and scoped to user, seat, and match. Every command is authorized against authenticated identity, room membership, seat, phase, allowed target set, revision, monotonic client sequence, and command ID.

`InternalMatchState` is never serialized. The server derives `PublicMatchView`, `PrivatePlayerView`, or restricted `SpectatorView`, then removes all other recipients from private event partitions before delivery. Active role mappings and unrevealed seeds are encrypted at rest with keys outside logs and database records.

Inputs receive size limits and strict schema validation. Accepted command IDs are idempotent; rate limits apply per session, connection, and command category. Origin checks, heartbeat expiry, bounded queues, and slow-client disconnection protect WebSockets. Chat moderation is a separate application boundary; unrestricted private live-match messages are prohibited.

Logs, traces, browser diagnostics, and crash reports redact roles, objectives, hidden cards, seeds, tokens, and raw private payloads. Database roles follow least privilege; secrets come from runtime configuration, never source. Dependency and container scanning run in CI.

The seed commitment is published before randomness is consumed. Role selection, event ordering, target selection, and Sigil shuffling use domain-separated streams. The seed is revealed only in `EPILOGUE` so auditors can verify the commitment and replay. Spectators receive only delayed, explicitly safe projections.
