# Development Workflow

1. Read the GDD, `AGENTS.md`, and the relevant architecture document.
2. Create a small feature-local change; update versioned contracts before or with wire changes.
3. Add unit tests for domain behavior and contract/integration tests at boundaries.
4. Run `make lint`, `make test`, and `make build`.
5. Apply SQL changes through reversible, numerically ordered migrations.
6. Add or supersede an ADR when changing an accepted architectural decision.

Frontend and backend release independently. Contract compatibility is additive within a protocol version; breaking changes require a new version and a compatibility window. Never edit generated artifacts manually. Reviews must call out hidden-information effects, failure recovery, performance impact, dependencies, and GDD deviations.
