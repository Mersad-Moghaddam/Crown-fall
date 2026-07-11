# ADR-0004: Realtime Protocol

**Status:** Accepted

## Context

Matches require bidirectional low-latency commands, ordered events, reconnect, and private payload filtering.

## Decision

Use versioned JSON envelopes over WebSocket, documented in AsyncAPI. Commands carry idempotency ID, expected revision, and monotonic sequence. Events carry authoritative revision/sequence and public/private partitions. HTTP remains for coarse resources.

## Alternatives considered

Polling cannot meet synchronized interaction needs. Server-sent events are one-way. A binary protocol reduces payload but makes early inspection/evolution harder.

## Consequences and follow-up

Strict runtime validation and compatibility policy are mandatory. Implement authentication, heartbeat, resumable event tails, payload compression, and generated contract conformance tests.
