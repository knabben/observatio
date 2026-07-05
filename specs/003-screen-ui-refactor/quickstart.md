# Quickstart: Screen Refactoring & UI Tech-Debt Remediation

How to run, test, and verify the refactored screens.

## Prerequisites

- Node + pnpm (frontend package manager per Constitution)
- Backend reachable (defaults to `http://localhost:8080`; override with `NEXT_PUBLIC_*` — see
  `contracts/environment-config.md`)

## Run

```bash
make run-frontend        # or: cd front && pnpm dev
make run-backend         # live data source
```

Production is the single binary — no separate frontend server, endpoints are same-origin:

```bash
make build            # builds frontend, embeds it, compiles output/observatio
./output/observatio   # serves UI + API/WS on one origin (any host/port, no rebuild)
make verify-binary    # build + launch + smoke-test the embed/serve end-to-end (US6)
```

Split dev mode only (frontend :3000 → backend :8080) uses the `NEXT_PUBLIC_*` overrides, which `run-frontend`
sets; they are unset in `make build`.

## Update dependencies (within-major, safe)

```bash
# frontend
cd front && pnpm update && pnpm add next@"^15" eslint-config-next@"^15"
# backend (safe modules only; k8s/CAPI/controller-runtime deferred)
cd webserver && go get github.com/gin-gonic/gin@latest github.com/gorilla/websocket@latest \
  github.com/spf13/cobra@latest github.com/stretchr/testify@latest && go mod tidy
```

Then confirm green: `make run-tests-frontend && make run-tests-backend && make build`.

## Test

```bash
make run-tests-frontend  # or: cd front && pnpm test
```

New/updated Jest tests (Testing Library + `MantineProvider` via `utils/test-render.tsx`) MUST cover, per
changed component:
- **Partial data**: resource missing `metadata`/`status`/`conditions`/`paused` renders without throwing.
- **Empty collection**: list renders a labeled empty state, not a header-only table.
- **Error/stuck**: socket connect-but-no-data resolves to empty/error within the bounded threshold; HTTP
  not-ok surfaces an error; empty frame does not clear a populated list.
- **Zero values**: numeric `0` renders as data; no stray `0` leaks into markup.
- **Status tri-state**: healthy / notready / unknown map to distinct, accessible indicators.

## Verify each screen (acceptance walkthrough)

| Screen | Check |
|--------|-------|
| Dashboard overview | Topology fits panel + `fitView`/controls; summary/versions/class tables show empty states; no horizontal scroll at laptop width |
| Clusters (+ infra) | No crash on partial cluster; spec panel full-width (no blank half); table scrolls in-container; keyboard-selectable rows |
| Machines (+ infra) | `0` cores renders cleanly; status unknown ≠ failed; `Age` labeled correctly |
| Machine Deployments | Unknown availability shown as unknown; stable row keys on re-sort |
| Shared shell | Nested route highlights parent nav; icon nav has accessible names; nav collapses on narrow viewport; "Search" filters; AI panel renders safely, stays in card, collapses |
| Single binary (US6) | `make verify-binary` passes; `./output/observatio` serves UI + API/WS same-origin; relocating host/port needs no rebuild; broken embed fails the target |

## Responsive check

Load every screen at ~1440px, ~1280px, ~768px: assert no horizontal page scroll, no permanently-empty
half-panels, columns stack, tables scroll within their containers.

## Definition of done

- All acceptance scenarios (spec US1–US6) pass.
- `make run-tests-frontend`, `make run-tests-backend`, `make build`, and `make verify-binary` succeed
  (Constitution: Development Workflow).
- Dependencies refreshed within-major (frontend + backend) with tests + build green.
- Success criteria SC-001…SC-011 verified.
