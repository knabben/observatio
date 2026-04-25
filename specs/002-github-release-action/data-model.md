# Data Model: Automated Release Publishing

**Feature**: `002-github-release-action`
**Date**: 2026-04-25

> This feature is CI/CD tooling; there is no persistent application data model.
> The entities below describe the conceptual model the workflow operates on.

## Entities

### Release Trigger

The git event that initiates the entire release pipeline.

| Attribute    | Type    | Description                                             |
|--------------|---------|---------------------------------------------------------|
| event_type   | string  | Always `push`                                          |
| ref          | string  | Full git ref, e.g. `refs/tags/v1.2.0`                 |
| ref_name     | string  | Short name extracted from ref, e.g. `v1.2.0`          |
| tag_pattern  | glob    | `v*` — only tags matching this pattern trigger the job |
| sha          | string  | Commit SHA the tag points to                           |

**Guard**: Pipeline MUST abort if `ref_name` does not start with `v`.

---

### Pipeline Run

A single execution of the release workflow, scoped to one trigger event.

| Attribute     | Type    | Description                                           |
|---------------|---------|-------------------------------------------------------|
| run_id        | string  | GitHub Actions run identifier                         |
| trigger       | Trigger | The event that initiated this run                     |
| status        | enum    | `in_progress` → `success` \| `failure`                |
| steps         | Step[]  | Ordered list of build and publish steps               |
| started_at    | datetime| Timestamp when the run began                          |

**State transitions**:
```
[queued] ──checkout──► [building] ──make build──► [built] ──publish──► [released]
                                                              │
                                                   [failure] (no release created)
```

---

### Release Asset

The compiled binary that is attached to the GitHub Release.

| Attribute      | Type   | Description                                             |
|----------------|--------|---------------------------------------------------------|
| source_path    | path   | `output/observatio` — produced by `make build`         |
| published_name | string | `observatio-<ref_name>-linux-amd64` (e.g. `observatio-v1.2.0-linux-amd64`) |
| platform       | string | `linux-amd64` (fixed for v1 scope)                     |
| size_bytes     | int    | File size of the compiled binary                       |
| executable     | bool   | Always true — binary is uploaded without wrapping      |

---

### GitHub Release

The release entry on the repository's Releases page.

| Attribute     | Type         | Description                                          |
|---------------|--------------|------------------------------------------------------|
| tag_name      | string       | Matches the trigger `ref_name`, e.g. `v1.2.0`       |
| name          | string       | Display name, defaults to tag name                   |
| body          | string       | Auto-generated release notes (GitHub default)        |
| draft         | bool         | Always false — releases are published immediately    |
| prerelease    | bool         | False for stable tags; determined by tag format      |
| assets        | Asset[]      | Exactly one asset per run (the production binary)    |
| created_by    | string       | `GITHUB_TOKEN` (automated, not a human actor)        |

---

### Workflow File

The declarative definition of the pipeline, stored in the repository.

| Attribute    | Type   | Description                                               |
|--------------|--------|-----------------------------------------------------------|
| path         | path   | `.github/workflows/release.yml`                          |
| trigger      | string | `on: push: tags: ['v*']`                                 |
| permissions  | map    | `contents: write` required for release creation          |
| runner       | string | `ubuntu-latest`                                          |
| dependencies | list   | `001-unified-build-script` features (`make build`) MUST exist |
