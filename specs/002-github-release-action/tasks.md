---
description: "Task list for Automated Release Publishing"
---

# Tasks: Automated Release Publishing

**Input**: Design documents from `specs/002-github-release-action/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Dependency note**: `001-unified-build-script` (make build) MUST be complete before this
feature can be fully validated end-to-end. Task T014 (`make build` step) depends on it.

**Organization**: Tasks grouped by user story. US2, US3, and US4 are structural properties
of the US1 `release.yml` implementation; they are covered by configuration choices and
verification tasks within their respective phases.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no shared dependencies)
- **[Story]**: Which user story this task belongs to (US1–US4)

---

## Phase 1: Setup

**Purpose**: Confirm baseline infrastructure before any changes.

- [x] T001 Verify `.github/workflows/` directory exists in project root

---

## Phase 2: Foundational — Fix Existing CI Workflow

**Purpose**: Repair `build.yml` to use `pnpm` and current action versions. This is a
prerequisite fix that MUST land before the release workflow, as it establishes the
correct tool setup pattern for the new workflow to follow.

**⚠️ CRITICAL**: All user story phases depend on `build.yml` being correct first.

- [x] T002 Upgrade `actions/checkout@v2` → `v4` in all three jobs in `.github/workflows/build.yml`
- [x] T003 [P] Upgrade `actions/setup-node@v2` → `v4` in `lint-test-frontend` and `build` jobs in `.github/workflows/build.yml`
- [x] T004 [P] Upgrade `actions/setup-go@v2` → `v5` with `go-version: '1.24.x'` in `lint-test-backend` and `build` jobs in `.github/workflows/build.yml`
- [x] T005 Add `pnpm/action-setup@v4` step (before the Node.js step) in the `lint-test-frontend` job in `.github/workflows/build.yml`
- [x] T006 Add `pnpm/action-setup@v4` step (before the Node.js step) in the `build` job in `.github/workflows/build.yml`
- [x] T007 Replace `npm install --include=dev` with `pnpm install --frozen-lockfile` in the `lint-test-frontend` job in `.github/workflows/build.yml`
- [x] T008 Replace `npm install` with `pnpm install --frozen-lockfile` in the `build` job in `.github/workflows/build.yml`

**Checkpoint**: `build.yml` uses pnpm throughout and all actions are at current major versions.

---

## Phase 3: User Story 1 — Tag-Triggered Release Pipeline (Priority: P1) 🎯 MVP

**Goal**: A `v*` tag push creates a GitHub Release with the compiled binary attached.

**Covers US3**: Tag-only trigger is implemented by the `on: push: tags: ['v*']` configuration
in T009. Branch pushes do not fire this workflow by design.

**Independent Test**: Push tag `v0.0.1-test` to a fork; verify GitHub Release is created
with asset `observatio-v0.0.1-test-linux-amd64` attached. Delete the test release/tag after.

### Implementation for User Story 1

- [x] T009 [US1] Create `.github/workflows/release.yml` with trigger `on: push: tags: ['v*']` and `permissions: contents: write`
- [x] T010 [P] [US1] Add `actions/checkout@v4` step (with `fetch-depth: 0`) to the release job in `.github/workflows/release.yml`
- [x] T011 [P] [US1] Add `actions/setup-go@v5` step with `go-version: '1.24.x'` to the release job in `.github/workflows/release.yml`
- [x] T012 [P] [US1] Add `actions/setup-node@v4` step with `node-version: '22'` and `pnpm/action-setup@v4` to the release job in `.github/workflows/release.yml`
- [x] T013 [US1] Add `pnpm install --frozen-lockfile` run step scoped to `front/` in `.github/workflows/release.yml`
- [x] T014 [US1] Add `make build` run step to produce `output/observatio` in `.github/workflows/release.yml`
- [x] T015 [US1] Add binary rename step `cp output/observatio observatio-${{ github.ref_name }}-linux-amd64` in `.github/workflows/release.yml`
- [x] T016 [US1] Add `softprops/action-gh-release@v2` step with `files: observatio-${{ github.ref_name }}-linux-amd64` in `.github/workflows/release.yml`

**Checkpoint**: Tag push produces GitHub Release with named binary asset. US3 validated
(non-tag pushes do not trigger the workflow).

---

## Phase 4: User Story 2 — Build Failure Blocks Release (Priority: P1)

**Goal**: Any build step failure prevents the release from being published.

**Independent Test**: Comment out the `run: make build` step body to simulate a build
failure; push a tag; verify no GitHub Release is created.

### Implementation for User Story 2

- [x] T017 [US2] Review step ordering in `.github/workflows/release.yml` — verify no step has `continue-on-error: true`; the publish step (T016) MUST appear after `make build` (T014) with no error bypass
- [x] T018 [US2] Confirm `softprops/action-gh-release@v2` step in `.github/workflows/release.yml` has no `if: always()` condition that would run even on prior step failure

**Checkpoint**: Simulated build failure leaves the Releases page unchanged.

---

## Phase 5: User Story 4 — Executable Release Asset (Priority: P2)

**Goal**: The published binary runs on a clean Linux machine with no extra installation.

**Covers US4**: Asset executability depends on the binary being statically linked
(`CGO_ENABLED=0`) which is already enforced in the Makefile.

**Independent Test**: Download the release binary; `chmod +x` it; run on a machine
without Go or Node installed; verify it starts and serves port 8080.

### Implementation for User Story 4

- [x] T019 [US4] Verify `CGO_ENABLED=0` is set in the `build` Makefile target in `Makefile` (static linking required for portable Linux binary)
- [x] T020 [P] [US4] Add a `chmod +x` step after the binary rename (T015) in `.github/workflows/release.yml` to ensure the uploaded asset has the executable bit set

**Checkpoint**: Downloaded binary is directly executable on a clean Linux host.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and verification across all user stories.

- [x] T021 [P] Update root `README.md` Releases/download section to reference the asset naming convention (`observatio-<version>-linux-amd64`)
- [x] T022 Run a full local `make build` to confirm `output/observatio` is produced before pushing the first real tag
- [x] T023 [P] Validate the final `.github/workflows/release.yml` YAML syntax with `yamllint` or the GitHub Actions schema validator

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Phase 1 — MUST complete before Phase 3+
- **US1/US3 (Phase 3)**: Depends on Phase 2 (correct pnpm setup establishes the pattern)
- **US2 (Phase 4)**: Depends on Phase 3 — T017/T018 review the release.yml from Phase 3
- **US4 (Phase 5)**: Depends on Phase 3 — T019/T020 extend the release.yml from Phase 3
- **Polish (Phase 6)**: Depends on all story phases complete

### User Story Dependencies

- **US1 (P1)**: Core workflow creation — all other stories depend on it
- **US2 (P1)**: Structural review of US1 output — can be done immediately after Phase 3
- **US3 (P2)**: Satisfied by US1 trigger config — no additional implementation
- **US4 (P2)**: Extends US1 with executable bit — can run in parallel with US2 review

### Parallel Opportunities

- T003, T004 can run in parallel (different jobs in build.yml)
- T010, T011, T012 can run in parallel (different steps added to release.yml)
- T019, T020 can run in parallel (different concerns within the executable asset story)
- T021, T023 can run in parallel (documentation and validation, no file conflicts)

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Fix build.yml
3. Complete Phase 3: Create release.yml
4. Complete Phase 4: Verify failure handling
5. **STOP and VALIDATE**: Push a test tag to confirm release is created
6. Merge if release publishes correctly

### Incremental Delivery

1. Fix build.yml → CI passes with pnpm ✅
2. Add release.yml → tag push produces a release ✅ (MVP)
3. Verify failure handling → quality gate confirmed ✅
4. Add executable bit → asset works on clean host ✅
5. Polish → docs updated, YAML validated ✅

---

## Notes

- [P] tasks = different files or non-conflicting sections, safe to parallelize
- US3 has no implementation tasks — it is the trigger configuration in T009
- `make build` (T014) depends on `001-unified-build-script` being complete
- Do not add `continue-on-error: true` to any step in the release workflow
- `secrets.GITHUB_TOKEN` is used implicitly by `softprops/action-gh-release@v2`; no
  additional secret configuration is required
