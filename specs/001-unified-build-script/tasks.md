---
description: "Task list for Unified Build System"
---

# Tasks: Unified Build System

**Input**: Design documents from `specs/001-unified-build-script/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, quickstart.md ✅

## Format: `[ID] [P?] [Story?] Description`

---

## Phase 1: Setup

- [X] T001 Verify/create `scripts/` directory at project root
- [X] T002 Verify `webserver/internal/web/handlers/build/` embed target directory exists

---

## Phase 2: Foundational — Prerequisite Validation Script

**Purpose**: Shared script called by every Makefile target to validate tools before work begins.

- [X] T003 Create `scripts/check-prereqs.sh` — validates go ≥1.24, node ≥22, pnpm presence with actionable error messages and per-scope filtering (--go, --node flags)

---

## Phase 3: User Story 1 — Validate Development Environment (Priority: P1) 🎯 MVP

**Goal**: `make check-prereqs` reports all tool issues before any build step runs.

**Independent Test**: Run `make check-prereqs` with Go uninstalled (or PATH removed); verify
non-zero exit and message naming the missing tool.

### Implementation for User Story 1

- [X] T004 [US1] Add `check-prereqs` target to `Makefile` that calls `scripts/check-prereqs.sh` for all required tools
- [X] T005 [US1] Make `scripts/check-prereqs.sh` executable (`chmod +x`)

**Checkpoint**: `make check-prereqs` exits 0 on valid setup, non-zero with tool-specific message on failure.

---

## Phase 4: User Stories 2 & 3 — Independent Backend and Frontend Builds (Priority: P1)

**Goal**: `make build-backend` and `make build-frontend` work independently without the other side present.

**Independent Test (backend)**: Remove `front/output/`; `make build-backend` still completes and produces a binary stub.
**Independent Test (frontend)**: `make build-frontend` runs with no Go toolchain on PATH and produces `front/output/`.

### Implementation

- [X] T006 [P] [US2] Add `build-backend` target to `Makefile` — runs go prereq check then `CGO_ENABLED=0 go build`
- [X] T007 [P] [US3] Add `build-frontend` target to `Makefile` — runs node/pnpm prereq check, `pnpm install --frozen-lockfile`, then `pnpm run build`

**Checkpoint**: Each target runs independently and exits non-zero on its own failure.

---

## Phase 5: User Story 4 — Unified Production Binary (Priority: P2)

**Goal**: `make build` produces a single binary with embedded frontend in one command.

**Independent Test**: Run `make build` on clean checkout; `output/observatio` exists and
starts serving port 8080 with `--dev=false`.

### Implementation

- [X] T008 [US4] Fix `build` target in `Makefile` — gate on `check-prereqs`, call `build-frontend`, copy assets cleanly, call `build-backend`
- [X] T009 [US4] Replace asset copy in `Makefile` with clean `rm -rf` + `cp -r` of `front/output/` into `webserver/internal/web/handlers/build/`

**Checkpoint**: `make build` produces `output/observatio`; binary serves dashboard.

---

## Phase 6: Polish

- [X] T010 Add `test` target to `Makefile` that runs both `run-tests-backend` and `run-tests-frontend`

---

## Dependencies & Execution Order

- Phase 1 → Phase 2 → Phase 3 → Phase 4 (T006/T007 parallel) → Phase 5 → Phase 6
- T006 and T007 can run in parallel (different Makefile targets, no file conflicts)
