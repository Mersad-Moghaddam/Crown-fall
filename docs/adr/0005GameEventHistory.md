# ADR-0005: Match Event History

**Status:** Accepted

## Context

Crownfall needs replay, Chronicle generation, recovery, fairness audit, debugging, and dispute evidence without adopting event sourcing for the whole product.

## Decision

Persist an append-only canonical event stream per match with strictly increasing sequence and explicit public/private partitions. Store compressed snapshots every 20 events and every phase boundary; recover from the newest valid snapshot plus tail.

## Alternatives considered

Snapshots alone lose audit history. Full product event sourcing is disproportionate. Redis-only history is not durable.

## Consequences and follow-up

Event schemas require versioning and secret retention/encryption policy. Implement transactional append plus snapshot, archive after the GDD's 90-day hot period, replay verification, and Chronicle template generation.
