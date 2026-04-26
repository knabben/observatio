# Contract: Release Workflow

**Type**: GitHub Actions workflow definition
**Path**: `.github/workflows/release.yml`
**Feature**: `002-github-release-action`

## Trigger Contract

```yaml
on:
  push:
    tags:
      - 'v*'
```

**Guarantee**: The workflow fires if and only if a tag whose name begins with `v` is pushed.
Any other event (branch push, PR, manual dispatch) MUST NOT trigger this workflow.

## Permissions Contract

```yaml
permissions:
  contents: write
```

**Requirement**: The `contents: write` permission is mandatory for `softprops/action-gh-release`
to create a release and upload assets. The `GITHUB_TOKEN` automatically carries this permission
when the above is set on the job.

## Environment Contract

| Tool         | Provided by              | Version        |
|--------------|--------------------------|----------------|
| Go           | `actions/setup-go@v5`    | `1.23.x`       |
| Node.js      | `actions/setup-node@v4`  | `20.x`         |
| pnpm         | `pnpm/action-setup@v4`   | `latest`       |
| make         | `ubuntu-latest` runner   | system default |

## Input Contract

| Input          | Source                    | Example          |
|----------------|---------------------------|------------------|
| Source code    | `actions/checkout@v4`     | full repo        |
| Version tag    | `github.ref_name`         | `v1.2.0`         |
| Publish token  | `secrets.GITHUB_TOKEN`    | built-in secret  |

## Output Contract

| Output           | Condition       | Description                                    |
|------------------|-----------------|------------------------------------------------|
| GitHub Release   | build success   | Release created at `github.com/<repo>/releases/tag/<tag>` |
| Release asset    | build success   | Binary `observatio-<tag>-linux-amd64` attached |
| Failed pipeline  | build failure   | No release created; CI run marked failed       |

## Asset Naming Contract

```
observatio-{tag}-linux-amd64
```

Where `{tag}` is the exact value of `github.ref_name` (e.g., `v1.2.0`).

**Examples**:
- Tag `v1.0.0` → asset `observatio-v1.0.0-linux-amd64`
- Tag `v2.3.1` → asset `observatio-v2.3.1-linux-amd64`

## Step Contract (ordered, each MUST succeed for next to run)

1. **checkout** — full history checkout (`fetch-depth: 0`)
2. **setup-go** — Go 1.23.x installed and on PATH
3. **setup-node** — Node 20.x installed
4. **setup-pnpm** — pnpm installed via `pnpm/action-setup@v4`
5. **install-deps** — `pnpm install --frozen-lockfile` in `front/`
6. **build** — `make build` exits 0; `output/observatio` exists
7. **rename-asset** — binary copied to `observatio-${{ github.ref_name }}-linux-amd64`
8. **publish-release** — `softprops/action-gh-release@v2` creates release + uploads asset

**Failure rule**: If any step exits non-zero, all subsequent steps are skipped and
the workflow run is marked as failed. No release is created or modified.
