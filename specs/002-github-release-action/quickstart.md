# Quickstart: Automated Release Publishing

**Feature**: `002-github-release-action`
**Date**: 2026-04-25

## Prerequisites

- Feature `001-unified-build-script` must be complete — `make build` must work locally.
- The repository must be hosted on GitHub with Actions enabled.
- No additional secrets need to be configured — `GITHUB_TOKEN` is provided automatically.

## How to Publish a Release

```bash
# 1. Ensure you are on main and everything is committed
git checkout main
git pull

# 2. Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

That's it. The release pipeline starts automatically within seconds of the tag push.

## Monitoring the Pipeline

1. Open the repository on GitHub
2. Go to **Actions** → **Release** workflow
3. Click the running workflow to see each step's output in real time

A successful run produces:
- A new entry on the **Releases** page
- A downloadable binary: `observatio-v1.0.0-linux-amd64`

## Downloading and Running the Release

```bash
# Download from GitHub Releases page
curl -L -o observatio \
  https://github.com/<owner>/<repo>/releases/download/v1.0.0/observatio-v1.0.0-linux-amd64

chmod +x observatio
./observatio serve
```

Open http://localhost:8080 — the full dashboard is served from this single binary.

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| Workflow does not trigger | Tag does not match `v*` | Use `v` prefix: `v1.0.0`, not `1.0.0` |
| `make build` step fails | Build system issue | Fix locally with `make build`, then push a new tag |
| Release created but asset missing | Upload step failed | Check workflow logs for asset path errors |
| `Resource not accessible by integration` | Wrong permissions | Ensure `permissions: contents: write` is in workflow YAML |

## Deleting a Bad Release

```bash
# Delete the release and its tag to retry
gh release delete v1.0.0 --yes
git tag -d v1.0.0
git push origin --delete v1.0.0

# Fix the issue, then re-tag
git tag v1.0.1
git push origin v1.0.1
```
