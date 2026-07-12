# System Context

The browser connects to the Go control plane over HTTP and WebSocket. The implemented slice owns one serialized in-memory room actor per match. PostgreSQL migrations and Redis configuration establish future durable-history and ephemeral-coordination boundaries, but neither is used by live match execution yet. Edge, CDN, persistence adapters, and OpenTelemetry export are planned boundaries.

Voice is a separate WebRTC/SFU plane. The browser may eventually use a server-authorized voice-room token, but audio never traverses the Go match server. This phase implements no voice signaling or media behavior.

See the [context diagram](diagrams/systemContext.mmd), [backend architecture](backendArchitecture.md), and [frontend architecture](frontendArchitecture.md).
