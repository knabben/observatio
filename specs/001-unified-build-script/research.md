# Research: Unified Build System

**Feature**: `001-unified-build-script`
**Branch**: `002-unified-build-script`
**Date**: 2026-04-25

## Findings

### Current Build Pipeline Audit

**Decision**: The existing `Makefile` build target is the reference implementation.
It is structurally sound but contains a critical bug and lacks prerequisite validation.

**Rationale**: Reading `Makefile`, `front/next.config.ts`, `front/package.json`,
and `webserver/internal/web/handlers/handlers.go` reveals the full pipeline:

```
front/                              webserver/
  pnpm run build                      //go:embed build/*
  → output/ (distDir in next.config)  ← handlers/build/
  → mv output/* handlers/build/       → CGO_ENABLED=0 go build
```

**Bug found**: The Makefile `build` target calls `npm run build` but `package.json`
scripts use `pnpm`. `pnpm` is the required package manager per the project constitution.
Using `npm` on a `pnpm`-lockfile project may produce inconsistent dependency resolution.

**Alternatives considered**:
- Replace Makefile with a standalone `scripts/build.sh` only → rejected; users already
  know `make build`, and Makefile is a standard Linux project interface.
- Keep Makefile as thin orchestrator, delegate stages to a shell script →
  **chosen approach**: simplest interface (`make build`) with clean stage isolation.

---

### Prerequisite Validation Strategy

**Decision**: Inline validation at the top of each Makefile target via a shared
`scripts/check-prereqs.sh` script.

**Rationale**: Running the check before any build step surfaces all problems
immediately (SC-002: within 5 seconds). The script checks:

| Tool    | Minimum version | Why required             |
|---------|-----------------|--------------------------|
| go      | 1.24            | Latest stable Go release  |
| node    | 22 (current LTS)| Latest Node.js LTS       |
| pnpm    | any             | Frontend package manager  |

**Detection approach**: `command -v <tool>` for presence; version string extraction
via `<tool> version` for version constraint checks using `sort -V` comparison.

**Alternatives considered**:
- Only check at full `make build` → rejected; independent backend/frontend targets
  also need their respective subsets of tools (FR-001 applies per stage).
- Docker-based build with all tools baked in → out of scope per constitution assumptions.

---

### Frontend Build Output and Asset Copy

**Decision**: Next.js static export outputs to `front/output/`; assets are copied
(not moved) to `webserver/internal/web/handlers/build/` before the Go embed step.

**Rationale**:
- `front/next.config.ts` sets `distDir: './output'` and `output: 'export'`, so
  the full static bundle lands in `front/output/`.
- The Go embed directive (`//go:embed build/*` in `handlers.go`) reads from
  `webserver/internal/web/handlers/build/`.
- The current Makefile uses `mv output/* ${BUILD_PATH}` which works but destroys
  the frontend output, making incremental rebuilds impossible. Using `cp -r` and
  then cleaning `build/` first is cleaner and safer.
- The existing `find ${BUILD_PATH} ! -name 'index.html' ! -name 'build' -type "f,d"`
  cleanup is fragile. A simple `rm -rf build/ && mkdir build/` before copying is
  equivalent and more readable.

**Alternatives considered**:
- Symlink `handlers/build` → `front/output` → rejected; `go:embed` does not follow
  symlinks outside the module root.

---

### Build Script Interface Design

**Decision**: Enhance the Makefile with clean independent targets. No separate
`scripts/build.sh` wrapper needed — the Makefile targets themselves are the script.

**Rationale**: The user requirement is "simplified experience like `make build`".
Introducing a separate shell script adds indirection. Instead:
- `make check-prereqs` — validates all tools
- `make build-frontend` — pnpm install + build only
- `make build-backend` — go build only (requires handlers/build/ to exist)
- `make build` — check-prereqs → build-frontend → build-backend in sequence
- `make run-backend` / `make run-frontend` — unchanged (development mode)
- `make test` — runs both backend and frontend test suites

Each target calls `check-prereqs` for its required tool subset before doing work.
Pipeline abort on first failure is enforced by Makefile's default error propagation.
