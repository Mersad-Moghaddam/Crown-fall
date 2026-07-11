# ADR-0002: Backend Architecture and Room Execution

**Status:** Accepted

## Context

Rules, hidden state, recovery, and replay need strong boundaries without premature distributed complexity. GDD §11.2 mandates one goroutine and serialized mailbox per match.

## Decision

Use a Go hexagonal modular monolith. Each room is owned by exactly one goroutine with a bounded mailbox. Domain/engine code remains infrastructure-free; future routing uses consistent hashing or a match directory.

## Alternatives considered

An actor framework adds supervision/distribution machinery before it is needed. A partitioned worker pool complicates room affinity and mutation ownership. Microservices add failure modes and secret-bearing network boundaries.

## Consequences and follow-up

Single-room mutation is race-free and testable, but blocking work must leave the actor through ports and return ordered results. Add mailbox saturation tests, durable actor recovery, and node-routing design before scaling horizontally.
