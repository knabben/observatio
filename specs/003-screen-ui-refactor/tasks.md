---
description: "Task list for Screen Refactoring & UI Tech-Debt Remediation"
---

# Tasks: Screen Refactoring & UI Tech-Debt Remediation

**Input**: Design documents from `/specs/003-screen-ui-refactor/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: INCLUDED — required by Constitution V (Test-Driven Quality) and mandated by spec/quickstart
(each changed component ships with Jest tests for partial-data, empty-collection, and error states).

**Organization**: Tasks are grouped by the six user stories (US1–US6) so each is independently
implementable and testable. Frontend lives in `front/app/`, backend in `webserver/`.

## Path Conventions

- Frontend: `front/app/**` (Next.js App Router; tests co-located as `*.test.tsx`, run via `pnpm test`)
- Backend: `webserver/**` (Go; tests as `*_test.go`, run via `go test ./...`)
- Test helper: `front/app/ui/dashboard/utils/test-render.tsx` (wraps `MantineProvider`)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Refresh dependencies within-major and establish a green baseline before any refactor (R11).

- [X] T001 [P] Refresh frontend deps within-major in `front/package.json`: `cd front && pnpm update && pnpm add next@"^15" eslint-config-next@"^15"`; keep lockfile (defer Mantine 9 / Jest 30)
- [X] T002 [P] Refresh backend deps within-major in `webserver/go.mod`: `cd webserver && go get github.com/gin-gonic/gin@latest github.com/gorilla/websocket@latest github.com/spf13/cobra@latest github.com/stretchr/testify@latest && go mod tidy` (defer k8s 0.36 / CAPI 1.13 / controller-runtime 0.24 / anthropic v1)
- [X] T003 Establish green baseline after updates: run `make run-tests-frontend`, `make run-tests-backend`, `make build`; fix only breakage introduced by the updates (no feature changes)

**Checkpoint**: Dependencies current within-major; build + tests green.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared building blocks, hardened data layer, and same-origin config that US1–US6 depend on.

**⚠️ CRITICAL**: No user-story work can begin until this phase is complete.

- [X] T004 [P] Create same-origin endpoint config in `front/app/lib/config.ts` (`API_URL` relative default, `WS_URL`/`WS_URL_CHATBOT` derived from `window.location`; `NEXT_PUBLIC_*` dev-only override) per `contracts/environment-config.md`
- [X] T005 Harden WebSocket client in `front/app/lib/websocket.tsx` (bounded reconnect: 8 attempts + exponential backoff to 30s, `onReconnectStop`→error; ignore empty/malformed frames — never clear list; import URLs from `config.ts`)
- [X] T006 Harden REST layer in `front/app/lib/data.tsx` (check `res.ok` before `json()`, wrap in try/catch, propagate errors; import `API_URL` from `config.ts`; remove hardcoded `URL` const)
- [X] T007 [P] Add `StatusState` type + `toStatusState()` tri-state helper (strict comparisons, absent→unknown) in `front/app/ui/dashboard/shared/status.ts`
- [X] T008 [P] Create `EmptyState` component in `front/app/ui/dashboard/shared/empty-state.tsx`
- [X] T009 [P] Create `ErrorState` component (message + optional keyboard-accessible retry) in `front/app/ui/dashboard/shared/error-state.tsx`
- [X] T010 [P] Create `StatusIndicator` (tri-state, no `processing` pulse on failed, accessible label) in `front/app/ui/dashboard/shared/status-indicator.tsx` (uses T007)
- [X] T011 Create `ChannelState` machine + `useResourceStream<T>` hook (connecting/ready/empty/error, 10s data-arrival timeout, `retry`) in `front/app/ui/dashboard/shared/resource-hooks.ts` (uses T005)
- [X] T012 [P] Create generic `ObjectTable<T>` (config-driven `ColumnDef`, stable `getRowKey`, `Table.ScrollContainer`, empty via `EmptyState`) in `front/app/ui/dashboard/shared/object-table.tsx` (uses T008)
- [X] T013 [P] Add shared `ColumnDef<T>` / `DetailFieldDef<T>` config types in `front/app/ui/dashboard/base/types.tsx`
- [X] T014 [P] Null-safe `FilterItems` (guard `metadata?.name`) in `front/app/dashboard/utils.tsx`
- [X] T015 Add sort-comparator guard for missing `metadata.name` in `front/app/ui/dashboard/base/lister.tsx`
- [X] T016 [P] Add root error boundary (`'use client'`) in `front/app/error.tsx`

**Checkpoint**: Shared primitives, hardened channels, and same-origin config ready — user stories can begin.

---

## Phase 3: User Story 1 - Screens never crash or hang on real cluster data (Priority: P1) 🎯 MVP

**Goal**: Every screen renders data / empty / error and never crashes or hangs on partial, empty, zero-value, or failed-connection data.

**Independent Test**: Point each screen at resources missing `status`/`metadata.name`/`conditions`/`paused`, at empty collections, and at a connect-but-no-data socket; confirm data, "no items", or an error message — never a crash or infinite spinner.

### Tests for User Story 1 ⚠️ (write first, ensure they fail)

- [X] T017 [P] [US1] Unit tests for `useResourceStream` states (connecting→ready/empty/error, 10s timeout, empty-frame no-op) in `front/app/ui/dashboard/shared/resource-hooks.test.tsx`
- [X] T018 [P] [US1] Tests for `EmptyState`/`ErrorState`/`StatusIndicator` + `toStatusState` (healthy/notready/unknown, zero-value, absent field) in `front/app/ui/dashboard/shared/status-indicator.test.tsx`
- [X] T019 [P] [US1] Partial-data + empty render tests for clusters (table, details, specification, infra) in `front/app/ui/dashboard/components/clusters/table.test.tsx`
- [X] T020 [P] [US1] Partial-data + empty + `0`-value render tests for machines (+infra, incl. `numCoresPerSocket=0` leak) in `front/app/ui/dashboard/components/machines/table.test.tsx`
- [X] T021 [P] [US1] Partial-data + empty + unknown-availability tests for mds in `front/app/ui/dashboard/components/mds/table.test.tsx`
- [X] T022 [P] [US1] Partial-data + empty tests for dashboard widgets (hierarchy/summary/versions/clusterclass) in `front/app/ui/dashboard/components/dashboard/dashboard.test.tsx`

### Implementation for User Story 1

- [X] T023 [P] [US1] Correct nullability (optional fields, `replicas: number`) in `front/app/ui/dashboard/components/clusters/types.tsx`
- [X] T024 [P] [US1] Correct nullability in `front/app/ui/dashboard/components/machines/types.tsx`
- [X] T025 [P] [US1] Correct nullability in `front/app/ui/dashboard/components/mds/types.tsx`
- [X] T026 [US1] Wire `BaseLister` to `useResourceStream` → `CenteredLoader`/`EmptyState`/`ErrorState(retry)` in `front/app/ui/dashboard/base/lister.tsx` (uses T011, T008, T009)
- [X] T027 [P] [US1] Null-safe access + empty state + `StatusIndicator` in `front/app/ui/dashboard/components/clusters/table.tsx` and `.../clusters/infra/infra-table.tsx`
- [X] T028 [P] [US1] Null-safe access (`paused`, arrays, empty `—`) in `front/app/ui/dashboard/components/clusters/details.tsx`, `.../specification.tsx`, and `.../clusters/infra/{infra-details,specification}.tsx`
- [X] T029 [P] [US1] Null-safe access + `StatusIndicator` + fix `0 && <JSX>` leak in `front/app/ui/dashboard/components/machines/table.tsx`, `.../machines/infra/infra-table.tsx`, `.../machines/infra/specification.tsx`
- [X] T030 [P] [US1] Null-safe access in `front/app/ui/dashboard/components/machines/details.tsx`, `.../specification.tsx`, `.../machines/infra/infra-details.tsx`
- [X] T031 [P] [US1] Null-safe access + `StatusIndicator` + strict comparisons (no `== 0`) in `front/app/ui/dashboard/components/mds/{table,details,specification}.tsx`
- [X] T032 [P] [US1] Guard `conditions.map`/`status`, add empty states, `generation.toString()`, remove `@ts-expect-error` in `front/app/ui/dashboard/components/dashboard/{clusterhierarchy,clustersummary,clusterversions,clusterclass}.tsx`
- [X] T033 [US1] Guard empty `tabs` + accurate rendering in `front/app/ui/dashboard/base/details.tsx`

**Checkpoint**: All screens survive partial/empty/zero/error data — MVP is independently shippable.

---

## Phase 4: User Story 2 - Layouts render correctly across sizes and data volumes (Priority: P1)

**Goal**: No horizontal page scroll, no half-empty panels, columns stack, tables scroll in-container, topology fits, at desktop/laptop/tablet.

**Independent Test**: Load each screen at ~1440/1280/768px with small and large data; verify no horizontal scroll, no permanently-empty half-panels, stacking, and in-container table scroll.

### Tests for User Story 2 ⚠️

- [X] T034 [P] [US2] Tests asserting `ObjectTable` renders a scroll container and detail panels use full width / responsive cols in `front/app/ui/dashboard/shared/object-table.test.tsx`

### Implementation for User Story 2

- [X] T035 [P] [US2] Responsive dashboard grid (`span={{base:12,md:...}}`) + topology container `width:100%` with `fitView`/`Background`/`Controls`/empty state in `front/app/dashboard/page.tsx` and `front/app/ui/dashboard/components/dashboard/clusterhierarchy.tsx`
- [X] T036 [P] [US2] Responsive summary/versions grids + fluid heights + versions scroll container in `front/app/ui/dashboard/components/dashboard/{clustersummary,clusterversions}.tsx`
- [X] T037 [US2] Add `Table.ScrollContainer` + responsive behavior in `front/app/ui/dashboard/shared/object-table.tsx`
- [X] T038 [P] [US2] Full-width single-child spec panels + responsive `SimpleGrid`/detail grids in `front/app/ui/dashboard/components/{clusters,machines,mds}/specification.tsx` and infra `specification.tsx`
- [X] T039 [P] [US2] Scroll containment at all breakpoints (not `md:`-only) in `front/app/dashboard/layout.tsx`
- [X] T040 [P] [US2] Vertically center loader (`min-h`) in `front/app/ui/dashboard/utils/loader.tsx`

**Checkpoint**: Screens are responsive and overflow-safe across supported widths.

---

## Phase 5: User Story 3 - Consistent, accessible navigation and status feedback (Priority: P2)

**Goal**: Nested-route active nav, keyboard/AT-operable controls, accessible names, three-state status legibility.

**Independent Test**: Visit a nested route (parent highlights); traverse all clickable names/links/icons by keyboard; inspect indicators for healthy vs not-ready vs unknown.

### Tests for User Story 3 ⚠️

- [X] T041 [P] [US3] Tests: active nav on nested route + `aria-current` + icon `aria-label` in `front/app/ui/dashboard/nav-links.test.tsx`
- [X] T042 [P] [US3] Test: `ObjectTable` selectable row is a keyboard-focusable, labeled control in `front/app/ui/dashboard/shared/object-table.a11y.test.tsx`

### Implementation for User Story 3

- [X] T043 [US3] Nested-route active state (`startsWith` + root guard) + `aria-current="page"` + icon `aria-label`s in `front/app/ui/dashboard/nav-links.tsx`
- [X] T044 [P] [US3] Responsive sidenav collapse (`Burger`/`Drawer`) + responsive logo `sizes`/alt in `front/app/ui/dashboard/sidenav.tsx`
- [X] T045 [US3] Keyboard-focusable, `aria`-labeled selectable rows (button/`UnstyledButton`) in `front/app/ui/dashboard/shared/object-table.tsx` (after T037)
- [X] T046 [P] [US3] Contrast-checked header/status colors + fix low-contrast title in `front/app/ui/dashboard/utils/header.tsx`

**Checkpoint**: Navigation, keyboard operability, and status legibility meet WCAG AA expectations.

---

## Phase 6: User Story 4 - Consolidated screen components and consistent theming (Priority: P2)

**Goal**: One shared, configurable implementation per concern (list/table/detail/spec); single theme/typography source.

**Independent Test**: Confirm each resource area is a config over shared blocks and a change to shared logic reflects everywhere; colors/fonts resolve from one source.

### Tests for User Story 4 ⚠️

- [X] T047 [P] [US4] Tests: per-resource `ColumnDef` configs render through shared `ObjectTable`; a shared-logic change reflects across resources in `front/app/ui/dashboard/shared/consolidation.test.tsx`

### Implementation for User Story 4

- [X] T048 [US4] Create Mantine theme (accent scale consolidating the scattered greens, status colors, font vars) in `front/app/styles/theme.ts`
- [X] T049 [US4] Apply `theme` + reconcile color scheme with system preference in `front/app/layout.tsx`; remove dead `--font-geist-*` vars and hardcoded `body` font in `front/app/globals.css` (after T048)
- [X] T050 [P] [US4] Fix font identifiers/weights (real fonts + bold weight) in `front/app/styles/fonts.ts`
- [X] T051 [P] [US4] Reduce clusters (+infra) to `ColumnDef`/`DetailFieldDef` over shared blocks in `front/app/ui/dashboard/components/clusters/{lister,table,details,specification}.tsx`
- [X] T052 [P] [US4] Reduce machines (+infra) to shared config in `front/app/ui/dashboard/components/machines/{lister,table,details,specification}.tsx`
- [X] T053 [P] [US4] Reduce mds to shared config in `front/app/ui/dashboard/components/mds/{lister,table,details,specification}.tsx`
- [X] T054 [US4] Dedupe fetch hooks → shared stream/fetch helper and fix copy-pasted error strings in `front/app/ui/dashboard/components/dashboard/*.tsx`
- [X] T055 [P] [US4] Remove dead Tailwind classes (`text-bold`, `text-medium`) + hardcoded hex; use theme tokens across `front/app/ui/dashboard/components/**`
- [X] T056 [P] [US4] Fix stale JSDoc and relabel `Created`→`Age` across listers/details in `front/app/ui/dashboard/components/**`

**Checkpoint**: Duplication collapsed to shared components; theming centralized; bugs can't drift per-copy.

---

## Phase 7: User Story 5 - Controls behave as labeled (Priority: P3)

**Goal**: "Search" filters, status chips aren't toggles, AI panel is safe, contained, and collapsible.

**Independent Test**: Type in Search → list filters; click a status chip → no toggle; AI content renders safely, stays in card, expands and collapses.

### Tests for User Story 5 ⚠️

- [X] T057 [P] [US5] Tests: search filters visible list; status chips are read-only; AI message renders as safe text and expand/collapse toggles in `front/app/ui/dashboard/base/ai-troubleshooting.test.tsx`

### Implementation for User Story 5

- [X] T058 [US5] Real client-side text filter wired to lister state in `front/app/ui/dashboard/search.tsx` (+ integrate in `front/app/ui/dashboard/base/lister.tsx`)
- [X] T059 [P] [US5] Replace toggle `Chip`s with read-only status display in `front/app/ui/dashboard/components/dashboard/clusterclass.tsx`
- [X] T060 [US5] AI panel fixes in `front/app/ui/dashboard/base/ai-troubleshooting.tsx`: safe plain-text render (remove `dangerouslySetInnerHTML`, valid nesting), remove `AppShell`-in-grid + contain to card, reversible expand, open-socket guard + reset loading, `uuid` id, functional state updates, `WS_URL_CHATBOT` from config

**Checkpoint**: Controls match their labels and the AI panel is safe and contained.

---

## Phase 8: User Story 6 - Runs as a single self-contained binary (Priority: P2)

**Goal**: One make target builds, embeds, and verifies a binary that serves UI + API/WS same-origin.

**Independent Test**: On a clean checkout run the build+verify target; the binary serves UI (200) and API/WS same-origin; relocating host/port needs no rebuild; a broken embed fails the target.

### Tests for User Story 6 ⚠️

- [X] T061 [P] [US6] Go test: embedded FS serves `/` (200, SPA HTML) and SPA fallback → `index.html` for unknown non-API routes in `webserver/internal/web/handlers/system/spa_test.go`

### Implementation for User Story 6

- [X] T062 [US6] Verify same-origin wiring end-to-end (no hardcoded `localhost`): `front/app/lib/data.tsx`, `front/app/lib/websocket.tsx`, `front/app/ui/dashboard/base/ai-troubleshooting.tsx` all use `config.ts`
- [X] T063 [US6] Confirm/implement SPA fallback (serve `index.html` for unknown non-API routes) in `webserver/internal/web/handlers/system/spa.go`
- [X] T064 [US6] Add `verify-binary` target to `Makefile` (depends on `build`; launch `output/observatio` on a test port → assert UI root 200 + SPA fallback + live API/WS same-origin → teardown → non-zero on failure) per `contracts/build-verification.md`
- [X] T065 [P] [US6] Document dev-only `NEXT_PUBLIC_*` vars in `README.md` / `.env.example`

**Checkpoint**: `make verify-binary` proves the single-binary embed+serve pipeline end-to-end.

---

## Phase 9: Polish & Cross-Cutting Concerns

- [X] T066 [P] Run `make lint-frontend` and `make lint-backend`; fix findings
- [X] T067 Run full gate: `make test` (backend+frontend) + `make build` + `make verify-binary`; ensure all green
- [X] T068 Execute `specs/003-screen-ui-refactor/quickstart.md` acceptance walkthrough (all screens, responsive check, single-binary)
- [X] T069 [P] Remove remaining dead code/unused fonts; final WCAG AA contrast audit across screens

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (P1)**: none — start immediately.
- **Foundational (P2)**: depends on Setup — BLOCKS all user stories.
- **US1 (P3)**: depends on Foundational. 🎯 MVP.
- **US2 (P4)**: depends on Foundational; T037 builds on `ObjectTable` (T012). Best after US1 (shares table/spec files) but independently testable.
- **US3 (P5)**: depends on Foundational; T045 depends on T037 (US2) — same file `object-table.tsx`.
- **US4 (P6)**: depends on Foundational; consolidation assumes shared blocks exist (T012, T037, T045). Heaviest file churn.
- **US5 (P7)**: depends on Foundational; independent of US2–US4.
- **US6 (P8)**: depends on Foundational (T004 config); T062 depends on config wiring. Backend-side largely independent.
- **Polish (P9)**: after all desired stories.

### Key cross-task dependencies

- `object-table.tsx`: T012 → T037 → T045 (sequential, same file).
- `theme.ts`/`layout.tsx`: T048 → T049.
- `base/lister.tsx`: T015 → T026 → T058 (sequential, same file).
- `ai-troubleshooting.tsx`: T057 (test) → T060 (impl); T059 separate file.

### Parallel Opportunities

- Setup: T001 ∥ T002.
- Foundational: T004, T007, T008, T009, T010, T012, T013, T014, T016 are ∥ (distinct files); T005/T006 use config, T011 uses T005, T015 edits lister.
- US1 tests T017–T022 all ∥; impl T023/T024/T025 ∥; T027–T032 ∥ (distinct resource files); T026/T033 edit base files.
- US4 T051 ∥ T052 ∥ T053 (distinct resource dirs).
- Stories US1, US5, US6 can largely proceed in parallel by different developers after Foundational; US2→US3→US4 share `object-table.tsx` and should serialize on it.

---

## Parallel Example: User Story 1

```bash
# Tests first (all parallel — distinct files):
Task: "T017 useResourceStream state tests"
Task: "T018 EmptyState/ErrorState/StatusIndicator tests"
Task: "T019 clusters partial/empty tests"
Task: "T020 machines partial/empty/zero tests"
Task: "T021 mds tests"
Task: "T022 dashboard widget tests"

# Then null-safe type fixes in parallel:
Task: "T023 clusters/types.tsx nullability"
Task: "T024 machines/types.tsx nullability"
Task: "T025 mds/types.tsx nullability"
```

---

## Implementation Strategy

### MVP First (User Story 1)

1. Phase 1 Setup → 2 Foundational → 3 US1.
2. **STOP and VALIDATE**: partial/empty/zero/error data across every screen (no crash, no infinite spinner).
3. Ship — this alone removes the most damaging user-facing defects.

### Incremental Delivery (recommended P1 launch = US1 + US2)

1. Foundation → US1 (robustness) → US2 (responsive) = a complete P1 increment.
2. US3 (a11y/nav) → US4 (consolidation/theming) → US6 (single-binary gate) → US5 (control polish).
3. Finish with Phase 9 (lint, full gate incl. `verify-binary`, quickstart walkthrough).

### Notes

- [P] = different files, no incomplete-task dependency.
- Write each story's tests first and confirm they fail before implementing (Constitution V).
- Commit after each task or logical group; keep `make build`, `make test`, and (from US6 on) `make verify-binary` green.
- Consolidation (US4) is where most individual findings collapse — do it after the shared blocks are proven by US1–US3 to avoid rework.
