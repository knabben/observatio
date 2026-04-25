# Quickstart: Unified Build System

**Feature**: `001-unified-build-script`
**Date**: 2026-04-25

## Prerequisites

Before running any build command, ensure the following tools are installed:

| Tool  | Minimum version | Install reference                         |
|-------|-----------------|-------------------------------------------|
| Go    | 1.24            | https://go.dev/doc/install               |
| Node  | 22 (LTS)        | https://nodejs.org or nvm                |
| pnpm  | latest          | `npm install -g pnpm`                    |
| make  | any             | Standard on Linux; `apt install make`    |

## Validate Your Environment

```bash
make check-prereqs
```

Expected output on a valid setup:
```
[check] go      1.23.x ✓
[check] node    18.x.x ✓
[check] pnpm    x.x.x  ✓
All prerequisites satisfied.
```

If a tool is missing, the output shows the tool name and an install hint, then exits non-zero.

## Build the Frontend Only

```bash
make build-frontend
```

Installs node dependencies (via pnpm) and produces the static export bundle at
`front/output/`. Does not require Go to be installed.

## Build the Backend Only

```bash
make build-backend
```

Compiles the Go server binary to `output/observatio`. Requires `front/output/` to have
been populated first (run `make build-frontend` or `make build` for a full build).

## Full Production Build (single binary)

```bash
make build
```

Runs all stages in order: prerequisite check → frontend build → asset copy → backend compile.
The result is a single self-contained binary at `output/observatio`.

Start it:
```bash
./output/observatio serve
```

Open http://localhost:8080 — the full dashboard UI and all API endpoints are served from
this single binary.

## Development Mode

Run frontend and backend independently in separate terminals:

**Terminal 1 — backend**:
```bash
make run-backend what=serve
```

**Terminal 2 — frontend**:
```bash
make run-frontend
```

Frontend dev server: http://localhost:3000
Backend API: http://localhost:8080

## Run Tests

```bash
make run-tests-backend     # Go test suite
make run-tests-frontend    # Jest test suite
```

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `pnpm: command not found` | pnpm not installed | `npm install -g pnpm` |
| `go: version too old` | Go < 1.23 installed | Update Go |
| Binary starts but UI is blank | Frontend assets not copied | Run `make build` (not `make build-backend` alone) |
| `pattern build/*: no matching files` | `handlers/build/` is empty | Run `make build-frontend` first |
