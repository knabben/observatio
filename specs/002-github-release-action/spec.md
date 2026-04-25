# Feature Specification: Automated Release Publishing

**Feature Branch**: `003-github-release-action`
**Created**: 2026-04-25
**Status**: Draft
**Input**: User description: "I want, after all these settings of build scripts, create a GitHub action to deploy the artifact automatically on their release. This will allow adding each tag a new artifact runs these scripts and publishes automatically."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Publish a Release by Pushing a Tag (Priority: P1)

A maintainer creates a new version tag (e.g., `v1.2.0`) and pushes it to the repository.
Without any further manual steps, a production binary is built and attached to a new
GitHub Release page for that tag — ready for users to download.

**Why this priority**: This is the core value of the feature. Every other story depends on
this pipeline working correctly.

**Independent Test**: Push a tag to a test repository branch; verify that a GitHub Release
is created with the binary attached and the release name matches the tag.

**Acceptance Scenarios**:

1. **Given** a maintainer pushes a tag matching the pattern `v*` (e.g., `v1.0.0`),
   **When** the automated pipeline runs, **Then** a GitHub Release is created for that tag
   with the compiled binary attached as a downloadable asset.
2. **Given** a tag push triggers the pipeline, **When** the build succeeds, **Then** the
   release asset filename includes the tag version and target platform (e.g.,
   `observatio-v1.0.0-linux-amd64`).
3. **Given** the pipeline produces the binary, **When** the GitHub Release is published,
   **Then** the release is visible on the repository's Releases page within 10 minutes of
   the tag being pushed.

---

### User Story 2 - Build Failure Blocks the Release (Priority: P1)

When the build pipeline fails for any reason — missing prerequisite, compilation error,
test failure — no release is published and the maintainer is notified of the failure with
enough context to diagnose and fix the problem.

**Why this priority**: Publishing a broken binary is worse than publishing nothing.
The pipeline MUST be a quality gate, not just a distribution mechanism.

**Independent Test**: Introduce a deliberate compile error; push a tag; verify that no
GitHub Release is created and the pipeline reports a failure.

**Acceptance Scenarios**:

1. **Given** the build step fails during a tag-triggered pipeline run, **When** the failure
   is detected, **Then** no GitHub Release is created or updated for that tag.
2. **Given** a pipeline failure occurs, **When** the maintainer views the CI run,
   **Then** the failed step is clearly identified with a non-zero exit code and log output.
3. **Given** a release was not published due to failure, **When** the maintainer fixes the
   issue and pushes a new tag (e.g., `v1.0.1`), **Then** the pipeline runs again and
   publishes the corrected release.

---

### User Story 3 - Pipeline Runs Only on Tags, Not on Every Commit (Priority: P2)

Ordinary commits and branch pushes do NOT trigger the release pipeline. Only tag pushes
matching the versioned tag pattern trigger a release build.

**Why this priority**: Running a full build and attempting to publish on every commit
would create noise, consume CI resources unnecessarily, and pollute the releases page.

**Independent Test**: Push a commit to `main` without a tag; verify that no release
pipeline run is triggered.

**Acceptance Scenarios**:

1. **Given** a commit is pushed to any branch without a version tag, **When** the CI system
   processes the event, **Then** the release pipeline does NOT run.
2. **Given** only a tag matching `v*` is pushed, **When** the CI system processes the event,
   **Then** the release pipeline runs exactly once.
3. **Given** a tag is pushed that does NOT match the `v*` pattern (e.g., a test tag),
   **When** the CI system processes the event, **Then** the release pipeline does NOT run.

---

### User Story 4 - Release Asset is Downloadable and Executable (Priority: P2)

A user downloads the binary from the GitHub Releases page and can run it directly on
a Linux machine without installing any additional dependencies.

**Why this priority**: The single-binary distribution model (Principle III) must hold end
to end. If the published asset requires external tools to run, the deployment goal fails.

**Independent Test**: Download the release asset from GitHub; run it on a clean Linux
machine with only the kubeconfig present; verify the dashboard is accessible on port 8080.

**Acceptance Scenarios**:

1. **Given** a release asset is downloaded from GitHub, **When** it is marked executable
   and run on a Linux machine, **Then** it starts and serves the dashboard without error.
2. **Given** the binary is started, **When** a valid kubeconfig is present at the default
   location, **Then** the dashboard connects to the cluster and displays health data.

---

### Edge Cases

- What happens if a release for the same tag already exists on GitHub? The pipeline MUST
  either update the existing release or fail with a clear message — it MUST NOT silently
  create a duplicate.
- What happens if CI secrets (e.g., publish token) are not configured? The pipeline MUST
  fail with an actionable error message identifying the missing secret, not a cryptic
  permissions error.
- What happens if the build takes longer than a reasonable timeout? The pipeline MUST have
  a configurable upper time limit and fail clearly if exceeded.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The release pipeline MUST trigger automatically when a tag matching the
  pattern `v*` is pushed to the repository.
- **FR-002**: The pipeline MUST run the full build process (prerequisite check → frontend
  build → backend compile) before publishing any release asset.
- **FR-003**: If any build step fails, the pipeline MUST stop immediately and MUST NOT
  publish a release.
- **FR-004**: On a successful build, the pipeline MUST create a GitHub Release for the
  pushed tag and attach the compiled binary as a downloadable asset.
- **FR-005**: The release asset filename MUST include the version tag and the target
  platform identifier.
- **FR-006**: The pipeline MUST NOT trigger on non-tag pushes (branches, commits).
- **FR-007**: Build and publish failures MUST surface as failed pipeline runs with
  identifiable failed steps and non-zero exit codes.
- **FR-008**: The published binary MUST be self-contained and runnable on Linux without
  additional runtime dependencies.

### Key Entities

- **Release Trigger**: The git event (tag push matching `v*`) that initiates the pipeline.
- **Release Asset**: The compiled, self-contained binary attached to the GitHub Release.
  Named with version and platform (e.g., `observatio-v1.0.0-linux-amd64`).
- **GitHub Release**: The published release entry on the repository's Releases page,
  created automatically by the pipeline for each successful tag build.
- **Pipeline Run**: A single execution of the automated build-and-publish workflow,
  scoped to one tag push event.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Pushing a version tag results in a published GitHub Release with a
  downloadable binary within 10 minutes, with no manual steps required from the maintainer.
- **SC-002**: A build failure on a tag push produces zero published release assets for
  that tag — the Releases page MUST NOT contain an empty or partial release.
- **SC-003**: A non-tag commit push triggers zero release pipeline executions.
- **SC-004**: The downloaded binary runs on a fresh Linux machine in under 5 seconds and
  serves the dashboard on port 8080 without requiring any additional installation steps.
- **SC-005**: A misconfigured secret or missing prerequisite causes the pipeline to fail
  with an error message that a maintainer can diagnose without reviewing raw CI logs.

## Assumptions

- The repository is hosted on GitHub and has access to GitHub Actions.
- The publish token/secret required to create GitHub Releases is stored as a repository
  secret and does not need to be provisioned as part of this feature.
- The target release platform is Linux amd64; multi-platform builds (arm64, darwin) are
  out of scope for this feature version.
- The build system defined in the Unified Build System feature (`001-unified-build-script`)
  is a prerequisite; this pipeline depends on `make build` being available and correct.
- Pre-release tags (e.g., `v1.0.0-rc1`) are treated the same as stable releases by this
  pipeline; distinguishing pre-release vs stable is out of scope.
- The CI environment provides all required tools (Go, Node.js, pnpm); the pipeline MUST
  validate prerequisites at runtime using the same check introduced in the build system feature.
