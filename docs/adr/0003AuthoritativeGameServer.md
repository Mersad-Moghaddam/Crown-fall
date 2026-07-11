# ADR-0003: Authoritative Game Server

**Status:** Accepted

## Context

Social deduction depends on secret integrity, simultaneous reveal, validated deadlines, and auditable randomness.

## Decision

Clients send intentions only. The server validates actions, advances the state machine, owns timers and randomness, selects endings, and builds explicit recipient projections. The engine is deterministic and infrastructure-free.

## Alternatives considered

Client authority is cheat-prone. Shared authority risks disagreement and leakage. Hiding a full state object in client UI is not security.

## Consequences and follow-up

All interactions require acknowledged state and latency-aware presentation. Add the complete GDD transition matrix, timeout policies, content validators, and security tests as gameplay is implemented.
