# Instructions for Coding Agents

1. Read `crownfallGameDesignDocument.md` before changing game behavior. Uppercase GDD state identifiers are normative. Use `EPILOGUE`, with the documented Figure 3 discrepancy; never create alternate canonical names.
2. Read the relevant architecture document. Preserve frontend/backend independence and exchange contracts only through versioned API schemas.
3. Use camelCase for project-authored multiword files and directories. Conventional exceptions include `README.md`, `AGENTS.md`, `LICENSE`, `Makefile`, `package.json`, `package-lock.json`, TypeScript tool configs, `go.mod`, `go.sum`, Go's required `_test.go` suffix (use names such as `engineTest_test.go`), `Dockerfile`, `.github/workflows`, dotfiles, PascalCase React/class files, and numeric migration prefixes. Multiword Go filenames deliberately use camelCase; do not convert them to snake_case.
4. Never trust the frontend, serialize `InternalMatchState`, expose secret state, bypass state-machine validation, or place game rules in HTTP/WebSocket handlers or repositories.
5. Voice media must never pass through the Go match server. The current voice directory is an extension point only.
6. Prefer small feature-local changes. Do not introduce microservices, generic utility dumping grounds, dependencies without justification, or speculative abstractions.
7. Add tests for every rule change. Update OpenAPI/AsyncAPI for contract changes and ADRs for architectural decisions. Do not edit generated files manually.
8. Run formatting, linting, tests, and builds. Report every changed file and unresolved risk honestly.
