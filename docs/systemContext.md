# System Context

The browser connects through an edge/gateway boundary to the Go control plane over HTTP and WebSocket. The modular match server owns one serialized room actor per match, persists event history and snapshots in PostgreSQL, and uses Redis only for ephemeral coordination. Assets come from a future CDN boundary. Telemetry is exported through OpenTelemetry-compatible interfaces.

Voice is a separate WebRTC/SFU plane. The browser may eventually use a server-authorized voice-room token, but audio never traverses the Go match server. This phase implements no voice signaling or media behavior.

See the [context diagram](diagrams/systemContext.mmd), [backend architecture](backendArchitecture.md), and [frontend architecture](frontendArchitecture.md).
