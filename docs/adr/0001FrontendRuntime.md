# ADR-0001: Frontend Runtime

**Status:** Accepted

## Context

Crownfall needs accessible application flows and a GPU-backed animated tabletop. GDD §11.1 explicitly selects React, TypeScript, and PixiJS v8.

## Decision

React owns document UI and PixiJS v8 owns the table through a narrow lifecycle/projection bridge. Vite builds the independently deployed frontend. PixiJS never owns authoritative state.

## Alternatives considered

Phaser offers a fuller game framework but conflicts with the stated renderer and adds systems unnecessary for a server-driven table. DOM-only rendering limits particle/cinematic performance. PixiJS is retained without substitution.

## Consequences and follow-up

The team must maintain semantic DOM equivalents and explicit bridge tests. Profile low-end tablets, formalize asset manifests, and validate performance budgets with representative art.
