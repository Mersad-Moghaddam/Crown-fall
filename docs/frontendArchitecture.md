# Frontend Architecture

## Goals and runtime boundaries

The browser must present a responsive, accessible, cinematic table while treating the Go server as the sole rules authority. React owns authentication, lobby, room settings, profile, chat, menus, results/Chronicle, reconnect, routing, and accessibility overlays. PixiJS v8 owns table, seats, cards, vote/Quest animation, scenario presentation, particles, transitions, and board interaction. WebAudio is isolated behind a future audio controller; WebRTC/SFU voice is a placeholder boundary only.

Node 24, TypeScript, Vite, React, PixiJS, Zustand, TanStack Query, React Router, Zod, Vitest, Testing Library, Playwright, ESLint, and Prettier form the toolchain.

## Structure and dependencies

`app` may import every frontend layer. `pages` may import features, entities, game, network, and shared. `features` may import entities, network, the game bridge, and shared. `entities` imports shared only. `game` imports protocol types and game-local modules. `network` imports protocol and shared validation. `shared` imports no page, feature, entity, or game code.

Feature directories cover authentication, lobby management, matchmaking, room management, role reveal, council, voting, quest, trial, ending, chat, reconnect, and the empty voice extension. `game` contains scenes, systems, objects, animation, audio, effects, camera, assets, and the React bridge. Multiword paths use camelCase; component/class files use PascalCase.

## State and data flow

- TanStack Query owns cacheable HTTP server state.
- One authoritative match store owns the latest acknowledged public snapshot and current player's private projection.
- A connection store owns socket status, sequence, heartbeat, retry, and reconnect token state.
- Feature-local Zustand or component state owns UI selections only.
- Game-local controllers own animation queues and audio state; frame-frequency data never enters React.

On a server event, the network layer validates the envelope with Zod, checks sequence continuity, applies acknowledged public/private projections once, and queues semantic animation instructions. PixiJS interpolates from the old visual state to the new projection without delaying authoritative state. Client intentions show non-secret pending feedback only; accepted/rejected events decide the outcome.

## Connection and recovery

Authenticate the WebSocket during upgrade or its first bounded message. Heartbeats detect stale connections. Reconnect uses exponential backoff with jitter and a short-lived token, sending last acknowledged sequence. The server returns an event tail or compressed recipient-safe snapshot plus tail. Hydration replaces the authoritative store atomically, cancels obsolete animations, and resumes from the server phase/deadline. A dedicated reconnect route explains prolonged failure.

## Presentation systems

The match page creates one PixiJS application and disposes it on unmount. Ordinary updates modify objects rather than reconstructing a scene. Assets use versioned manifests, texture atlases, scenario bundles, integrity hashes, preload budgets, and deterministic fallbacks. Sound categories—music, ambience, effects, interface, voice—have independent gain; user gesture unlock and reduced-sensory settings are respected.

Layouts target desktop and responsive tablet, use logical CSS properties for RTL, and keep controls outside canvas where possible. Canvas interactions have semantic keyboard/DOM equivalents, visible focus, non-color cues, screen-reader announcements, captions/text alternatives, scalable text, and reduced-motion/effects modes.

## Errors, security, and performance

Protocol failure is fatal to the connection, not silently coerced. Recoverable network, asset, render, and application-boundary errors have distinct messages and correlation IDs. Browser logs and telemetry redact private payloads. The client never receives other players' secrets and never persists unrevealed private state beyond the session.

Targets requiring validation: compressed shell plus current scenario ≤3.5 MB; interactive lobby p75 ≤2.5 seconds; board 60 FPS with 30 FPS reduced-effects minimum; typical public event ≤4 KB and compressed snapshot ≤64 KB. Routes/assets are lazy, React subscriptions are narrow, network updates are separated from frames, and scene rebuilds are avoided.

## Testing and example lifecycle

Test components, hooks, stores, parsing, bridge lifecycle, accessibility, reconnection, visual smoke, and E2E journeys. Example: player submits `CAST_VOTE`; UI shows pending; server validates and emits an accepted event; the authoritative store updates; the bridge queues a simultaneous reveal; rejection restores controls with the server reason. PixiJS never decides whether the vote is legal.
