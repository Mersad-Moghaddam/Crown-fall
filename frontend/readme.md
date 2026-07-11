# Crownfall Frontend

The frontend owns navigation, accessible UI, lobby/room flows, connection lifecycle, and presentation. React owns application UI; PixiJS v8 owns the animated table. Neither validates rules or advances authoritative state.

## Stack and structure

Node 24, TypeScript, React, Vite, PixiJS, Zustand, TanStack Query, React Router, Zod, Vitest, Testing Library, Playwright, ESLint, and Prettier. Code is feature-oriented under `src/app`, `pages`, `features`, `entities`, `game`, `network`, and `shared`. `features/voice` is an unimplemented boundary.

State is separated into authoritative remote snapshots, private player projection, connection, local UI, animation, and audio state. Authoritative state is stored once; PixiJS receives render projections through `game/bridge`.

## Environment and commands

Copy `.env.example` to `.env`. `VITE_API_URL` selects HTTP and `VITE_WEBSOCKET_URL` selects the control socket.

```sh
npm ci
npm run dev
npm run typecheck
npm run lint
npm run formatCheck
npm test
npm run build
npm run testE2e
```

## Development rules

- Add a feature inside `features/<featureName>` and expose the smallest API needed by pages.
- Add a PixiJS scene under `game/scenes`; mount it only through a bridge that creates and destroys the Pixi lifecycle.
- Add WebSocket handlers under `network/websocket`, validate envelopes with Zod, update authoritative state before animations, and never infer command acceptance.
- Lazy-load non-match routes and scenario assets. Use atlases, avoid high-frequency React state, preserve 60 FPS with a 30 FPS reduced-effects floor, and keep shell plus current scenario at or below 3.5 MB compressed.
- Use semantic HTML outside the canvas, keyboard equivalents, visible focus, reduced motion, color-independent cues, and logical CSS properties for RTL readiness.

If the canvas is blank, check WebGL support and its host size. If protocol parsing fails, compare `backend/api/asyncapi.yaml`. If E2E browsers are absent, run `npx playwright install`.
