# Implementation Plan: Unified Build System

**Branch**: `002-unified-build-script` | **Date**: 2026-04-25 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/001-unified-build-script/spec.md`

## Summary

Repair and extend the existing `Makefile` build pipeline to deliver a simple,
reliable developer experience: `make build` produces a single Go binary with embedded
frontend, while `make check-prereqs`, `make build-frontend`, and `make build-backend`
enable independent stage execution. A shared prerequisite validation script is introduced
to catch missing or outdated tools before any build step runs.

## Technical Context

**Language/Version**: Go 1.24 · Bash (GNU)
**Primary Dependencies**: pnpm (frontend), cobra (backend CLI), gorilla/mux (router)
**Storage**: N/A (build tooling only)
**Testing**: `go test ./...` (backend) · Jest via `pnpm run test` (frontend)
**Target Platform**: Linux (host execution — no Docker)
**Project Type**: Build tooling / shell scripts + existing web service
**Performance Goals**: Full build < 5 min · Prereq check < 5 sec
**Constraints**: Must keep `make build` as the primary interface; pnpm required (not npm);
  minimum Go 1.24 and Node 22 (LTS)
**Scale/Scope**: Single developer machine; single project repo

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Observability-First | ✅ Pass | Build stages MUST log each step and surface errors with file/line context |
| II. Real-Time Visibility | ✅ N/A | Build tooling; streaming stdout satisfies spirit of principle |
| III. ClusterAPI Resource Model | ✅ N/A | No domain model changes; build artefacts only |
| IV. AI-Augmented Troubleshooting | ✅ N/A | No runtime AI integration in build scripts |
| V. Test-Driven Quality | ✅ Pass | Existing test suites (`go test`, Jest) MUST pass; `make build` MUST run them cleanly |

**Post-design re-check (Phase 1)**: No violations introduced. The new Makefile targets
do not alter source code, only the build pipeline. Test suites remain unchanged.

## Project Structure

### Documentation (this feature)

```text
specs/001-unified-build-script/
├── plan.md           # This file
├── research.md       # Phase 0 output
├── data-model.md     # Phase 1 output
├── quickstart.md     # Phase 1 output
└── tasks.md          # Phase 2 output (/speckit-tasks command)
```

### Source Code (repository root)

```text
Makefile                                        ← modified: new targets + pnpm fix
scripts/
└── check-prereqs.sh                            ← new: shared prereq validation script

webserver/internal/web/handlers/
└── build/                                      ← existing embed target dir (unchanged)
```

No new source packages. No structural changes to `webserver/` or `front/`.

## Complexity Tracking

> No constitution violations; no complexity justification required.
