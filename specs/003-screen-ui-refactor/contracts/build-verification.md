# Contract: Single-Binary Build & Verification

Covers US6 / FR-037–FR-039. The stack ships as one self-contained binary (`output/observatio`) with the
exported frontend embedded via `//go:embed build/*` and served same-origin with the API/WebSocket.

## Existing pipeline (unchanged)

```text
make build:
  build-frontend         # pnpm install --frozen-lockfile && pnpm run build  → front/output
  copy embed assets      # rm -rf webserver/internal/web/handlers/build; cp -r front/output/. → build/
  build-backend          # CGO_ENABLED=0 go build -o output/observatio .   (embeds build/*)
```

## New target: `make verify-binary`

**Contract** — the target MUST:
1. Depend on / run `make build` (clean-checkout runnable; no separate frontend dev server — FR-039).
2. Launch `output/observatio` as a background process on a test port.
3. Assert, over the same origin:
   - `GET /` → **HTTP 200** and returns the SPA HTML shell (embedded UI root). (FR-038)
   - An unknown non-API client route → **HTTP 200** serving `index.html` (SPA fallback intact).
   - A live API/WebSocket endpoint responds (e.g. WS upgrade on `/ws` succeeds, or a known REST route
     returns a non-5xx). (FR-038)
4. Terminate the launched process (always, even on failure).
5. Exit **non-zero** if any assertion fails — a missing/stale embed or broken SPA fallback fails the target,
   not a green build. (FR-037, edge case: missing/stale embed)

**Inputs**: none beyond a clean checkout + toolchain (`go`, `node`/`pnpm`).
**Outputs**: exit code (0 = seamless embed+serve verified); human-readable pass/fail log lines.

## Optional complement (backend unit test)

A Go test using `net/http/httptest` against the embedded `fs.FS` MAY assert `GET /` and the SPA fallback
resolve from the embed filesystem — fast, runs in `make run-tests-backend`. This complements, but does not
replace, `verify-binary` (which exercises the actual built binary end-to-end).

## Same-origin guarantee

Because the binary serves both the UI and API/WS, the frontend addresses them relatively (see
`environment-config.md`). `verify-binary` implicitly proves same-origin: the smoke checks hit the UI and an
API/WS endpoint on the **one** origin the binary listens on, with `NEXT_PUBLIC_*` unset. (SC-009, SC-011)
