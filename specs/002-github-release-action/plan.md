# Implementation Plan: Automated Release Publishing

**Branch**: `003-github-release-action` | **Date**: 2026-04-25 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/002-github-release-action/spec.md`

## Summary

Create `.github/workflows/release.yml` — a GitHub Actions workflow that triggers on any
`v*` tag push, runs the full build pipeline (`make build`), and publishes a self-contained
Linux binary as a GitHub Release asset. As a prerequisite fix, update the existing
`build.yml` to use `pnpm` (not `npm`) and upgrade action versions to current majors.

## Technical Context

**Language/Version**: Go 1.23.1 · Bash · YAML (GitHub Actions)
**Primary Dependencies**: `softprops/action-gh-release@v2` · `pnpm/action-setup@v4`
**Storage**: N/A — workflow is stateless; release assets stored by GitHub
**Testing**: Pipeline validated by pushing a test tag to a fork/branch
**Target Platform**: `ubuntu-latest` GitHub Actions runner (Linux amd64)
**Project Type**: CI/CD workflow (YAML) + fix to existing workflow
**Performance Goals**: Full release pipeline completes within 10 minutes of tag push
**Constraints**: Must use `GITHUB_TOKEN` only (no PAT); must use `pnpm` not `npm`;
  depends on `make build` from feature `001-unified-build-script`
**Scale/Scope**: One workflow, one job, eight sequential steps

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Observability-First | ✅ Pass | Each workflow step is named and logs to GitHub Actions; failures surface with step-level exit codes |
| II. Real-Time Visibility | ✅ N/A | CI/CD tooling; pipeline logs satisfy spirit of principle |
| III. ClusterAPI Resource Model | ✅ N/A | No domain model changes |
| IV. AI-Augmented Troubleshooting | ✅ N/A | No runtime AI integration |
| V. Test-Driven Quality | ✅ Pass | `make build` is a prerequisite; the existing test suite runs in `build.yml` before any release; release workflow itself is not bypassing tests |

**Post-design re-check (Phase 1)**: No violations. The release workflow calls `make build`
which compiles tested code. No production code changes.

## Project Structure

### Documentation (this feature)

```text
specs/002-github-release-action/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── workflow-contract.md   # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit-tasks command)
```

### Source Code (repository root)

```text
.github/workflows/
├── build.yml            ← modified: npm → pnpm, action version upgrades
└── release.yml          ← new: tag-triggered release pipeline
```

No changes to `webserver/`, `front/`, `Makefile`, or any application source.

## Complexity Tracking

> No constitution violations; no complexity justification required.

## Key Implementation Notes

### `release.yml` structure

```yaml
name: Release
on:
  push:
    tags: ['v*']

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23.x' }
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - uses: pnpm/action-setup@v4
        with: { version: latest }
      - run: cd front && pnpm install --frozen-lockfile
      - run: make build
      - run: cp output/observatio observatio-${{ github.ref_name }}-linux-amd64
      - uses: softprops/action-gh-release@v2
        with:
          files: observatio-${{ github.ref_name }}-linux-amd64
```

### `build.yml` fixes (same PR)

- Replace `npm install` → `pnpm install --frozen-lockfile` in all jobs
- Add `pnpm/action-setup@v4` step before Node.js usage
- Upgrade `actions/checkout@v2` → `v4`
- Upgrade `actions/setup-node@v2` → `v4`
- Upgrade `actions/setup-go@v2` → `v5`
