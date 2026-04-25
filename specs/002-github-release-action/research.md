# Research: Automated Release Publishing

**Feature**: `002-github-release-action`
**Branch**: `003-github-release-action`
**Date**: 2026-04-25

## Findings

### GitHub Actions Trigger for Tag Pushes

**Decision**: Use `on: push: tags: ['v*']` as the sole trigger for the release workflow.

**Rationale**: This is the canonical GitHub Actions pattern for version-tag-triggered
workflows. It fires exactly once per matching tag push and does not fire on branch pushes
or commit events — satisfying FR-006 and US3 precisely.

**Alternatives considered**:
- `on: release: types: [created]` — requires the maintainer to manually create a Draft
  Release on GitHub first; rejected because it requires manual steps (violates SC-001).
- `on: push:` without tag filter — fires on every commit; rejected (US3 explicitly prohibits this).

---

### Release Publishing Action

**Decision**: Use `softprops/action-gh-release@v2` for creating the GitHub Release and
uploading the binary asset in a single step.

**Rationale**: This action is the most widely adopted release-publishing action in the
GitHub ecosystem. It handles idempotent release creation (updates existing release for the
same tag rather than failing or duplicating), accepts a glob for asset paths, and uses
the built-in `GITHUB_TOKEN` — no extra secrets needed beyond what GitHub provides.

**Token**: `GITHUB_TOKEN` (built-in) with `permissions: contents: write` on the job is
sufficient to create releases and upload assets on the same repository.

**Alternatives considered**:
- `gh release create` via the GitHub CLI — equally valid but requires more scripting to
  handle the asset rename and idempotency; `softprops/action-gh-release` is simpler.
- `actions/create-release` + `actions/upload-release-asset` — deprecated by GitHub; rejected.

---

### Existing Workflow Bugs

**Decision**: Fix the existing `build.yml` as part of this feature — use `pnpm` instead
of `npm`, and upgrade action versions to current majors.

**Rationale**: Reading `.github/workflows/build.yml` reveals:
- `npm install` is used in all jobs; the project's `package.json` and constitution require `pnpm`.
- `actions/checkout@v2`, `actions/setup-node@v2`, `actions/setup-go@v2` are outdated
  (current: v4, v4, v5 respectively). Outdated actions receive no security patches.
- The `release.yml` workflow will introduce `pnpm` setup correctly; `build.yml` must
  match to avoid inconsistency.

**Fix scope**: Update `build.yml` in the same PR: switch to `pnpm/action-setup@v4` for
pnpm installation, update action versions, replace `npm install` → pnpm.

---

### Binary Naming and Platform Identification

**Decision**: Name the release asset `observatio-<TAG>-linux-amd64` where `<TAG>` is
extracted from `${{ github.ref_name }}` (e.g., `v1.2.0`).

**Rationale**: This naming convention is standard across Go CLI projects (e.g., kubectl,
helm, clusterctl). It allows users to identify version and platform at a glance without
unpacking an archive. No tar/zip wrapping for v1 — the binary is directly downloadable
(the project is a single statically linked executable).

**Alternatives considered**:
- Wrapping in `.tar.gz` — conventional but adds friction for a single binary. Out of scope v1.
- Platform matrix (linux/arm64, darwin) — out of scope per spec assumptions.

---

### Build Job Structure

**Decision**: The release workflow is a single job with sequential steps; it does NOT
reuse the CI `build` job via `needs`. It sets up the full tool chain from scratch.

**Rationale**: Release builds should be hermetic — a standalone, reproducible run that
does not depend on the state of parallel CI jobs. Reusing the `build` job's artefacts
via `actions/upload-artifact` would introduce coupling and cache invalidation complexity.

**Steps**:
1. Checkout (with full history for tag metadata)
2. Setup Go 1.23
3. Setup Node 20 + pnpm via `pnpm/action-setup@v4`
4. Install frontend deps (`pnpm install --frozen-lockfile`)
5. Run `make build` (produces `output/observatio`)
6. Rename binary to `observatio-${{ github.ref_name }}-linux-amd64`
7. Create GitHub Release + upload renamed binary via `softprops/action-gh-release@v2`
