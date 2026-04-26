# Data Model: Unified Build System

**Feature**: `001-unified-build-script`
**Date**: 2026-04-25

> This feature is build tooling; there is no persistent data store or runtime data model.
> The entities below describe the conceptual model that the build scripts operate on.

## Entities

### Build Stage

Represents a discrete, ordered step in the build pipeline. Each stage MUST be independently
runnable and independently failable.

| Attribute       | Type    | Description                                              |
|-----------------|---------|----------------------------------------------------------|
| name            | string  | Human-readable label (e.g., "check-prereqs")             |
| make_target     | string  | The `make <target>` invocation that runs this stage      |
| prerequisite_of | Stage[] | Stages that MUST complete before this stage runs         |
| tool_deps       | Tool[]  | Prerequisites required for this stage to execute        |
| exit_on_failure | bool    | Always true — no stage continues after a failure         |

**Stages (ordered)**:
1. `check-prereqs` — validates all required tools for the requested build scope
2. `build-frontend` — installs node deps and produces the static export bundle
3. `copy-assets` — clears `handlers/build/` and copies the frontend bundle into it
4. `build-backend` — compiles the Go binary with embedded assets

### Prerequisite (Tool)

A tool that MUST be present and at or above a minimum version before a build stage runs.

| Attribute       | Type   | Description                                         |
|-----------------|--------|-----------------------------------------------------|
| name            | string | Executable name (e.g., `go`, `pnpm`, `node`)        |
| min_version     | string | Minimum acceptable semver string (e.g., `1.23.0`)  |
| used_by         | Stage[]| Which build stages require this tool                |
| install_hint    | string | Short guidance shown when tool is missing            |

**Defined prerequisites**:
| Tool  | Min version | Used by stages                              |
|-------|-------------|---------------------------------------------|
| go    | 1.24.0      | check-prereqs, build-backend                |
| node  | 22.0.0      | check-prereqs, build-frontend               |
| pnpm  | any         | check-prereqs, build-frontend               |

### Embedded Asset Bundle

The compiled frontend output that is embedded inside the Go binary.

| Attribute      | Type   | Description                                              |
|----------------|--------|----------------------------------------------------------|
| source_dir     | path   | `front/output/` — Next.js static export output          |
| dest_dir       | path   | `webserver/internal/web/handlers/build/`                |
| embed_directive| string | `//go:embed build/*` in `handlers.go`                   |
| staleness      | bool   | True if source was rebuilt after last copy               |

**State transitions**:
```
[missing] ──build-frontend──► [built at source_dir]
                                       │
                              ──copy-assets──►  [copied to dest_dir]
                                                        │
                                           ──build-backend──► [embedded in binary]
```

### Production Binary

The final output artefact: a self-contained executable embedding the full frontend bundle.

| Attribute      | Type   | Description                                          |
|----------------|--------|------------------------------------------------------|
| output_path    | path   | `output/observatio` (relative to project root)       |
| embed_includes | path[] | All files under `webserver/internal/web/handlers/build/` |
| dev_mode_flag  | bool   | `--dev` flag disables static hosting at runtime      |
| cgo_enabled    | bool   | Always false (`CGO_ENABLED=0`) for portability       |
