# Feature Specification: Unified Build System

**Feature Branch**: `002-unified-build-script`
**Created**: 2026-04-25
**Status**: Draft
**Input**: User description: "Review all the base build and implementation of the projects, create a shell script for all the requirements: independence in both sides; check and understand all the internals; hookings to create a only a golang binary with front-end built in."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Validate Development Environment (Priority: P1)

A developer setting up the project for the first time runs a single prerequisite-check
command that inspects the local environment and reports exactly which tools are present,
which are missing, and what version constraints are violated — before any build step is
attempted.

**Why this priority**: Without environment validation, developers encounter cryptic mid-build
failures. Surfacing all problems upfront saves significant setup time and reduces support burden.

**Independent Test**: Run the prerequisite check in an environment with one tool deliberately
missing; the output MUST list the missing tool and stop before any build artefact is produced.

**Acceptance Scenarios**:

1. **Given** all required tools are installed at acceptable versions, **When** the prerequisite
   check runs, **Then** it exits successfully with a confirmation message and no warnings.
2. **Given** a required tool is absent, **When** the prerequisite check runs, **Then** it
   prints the name of the missing tool, explains where to obtain it, and exits with a non-zero
   code.
3. **Given** a tool is installed but below the minimum required version, **When** the check
   runs, **Then** it reports the installed version, the required version, and exits non-zero.

---

### User Story 2 - Build Backend Independently (Priority: P1)

A backend developer runs a single command to compile and verify the server component
without any dependency on the frontend source code or frontend build artefacts.

**Why this priority**: Backend and frontend developers work in parallel. Coupling their build
steps creates unnecessary bottlenecks and hides where failures originate.

**Independent Test**: Remove all frontend assets from the workspace; the backend build command
MUST still complete successfully and produce a runnable server binary.

**Acceptance Scenarios**:

1. **Given** the backend source is present and prerequisites pass, **When** the backend build
   command runs, **Then** a runnable server binary is produced with no frontend assets required.
2. **Given** a backend source file contains a compile error, **When** the build runs, **Then**
   the error is reported with file name and line number, and no partial binary is left behind.
3. **Given** the backend build succeeds, **When** the binary is started in development mode,
   **Then** it serves API endpoints on the configured port.

---

### User Story 3 - Build Frontend Independently (Priority: P1)

A frontend developer runs a single command to install dependencies and produce a production
bundle of the UI without requiring the backend to be compiled or running.

**Why this priority**: Same rationale as Story 2 — independent build cycles accelerate
iteration and isolate failures.

**Independent Test**: Run the frontend build command on a machine with no Go toolchain
installed; it MUST complete and produce a deployable UI bundle.

**Acceptance Scenarios**:

1. **Given** frontend source is present and node package manager is available, **When** the
   frontend build command runs, **Then** a deployable static bundle is produced.
2. **Given** a missing or incompatible node dependency, **When** the build runs, **Then** the
   error clearly identifies the package conflict and suggests remediation.
3. **Given** the frontend bundle is produced, **When** it is served against the backend API,
   **Then** all dashboard pages load without errors.

---

### User Story 4 - Produce a Unified Production Binary (Priority: P2)

A release engineer runs a single full-build command that validates prerequisites, builds the
frontend, embeds the frontend bundle inside the server binary, and outputs a single
self-contained executable that serves both API and UI — no separate static file deployment
required.

**Why this priority**: The single-binary deployment model is a key operational advantage of
this project. It simplifies distribution, container images, and deployment scripts.

**Independent Test**: Run the full-build command on a clean checkout; the resulting binary MUST
start and serve the dashboard UI and all API endpoints from port 8080 without any additional
files present.

**Acceptance Scenarios**:

1. **Given** prerequisites pass and all source is present, **When** the full build command runs,
   **Then** a single binary is produced that contains the embedded frontend and the full API.
2. **Given** the production binary is started with no external files, **When** a browser
   navigates to port 8080, **Then** the full dashboard UI loads and all API calls succeed.
3. **Given** a previous build artefact exists, **When** the full build runs, **Then** stale
   artefacts are cleaned before the new build begins.
4. **Given** the frontend build step fails mid-way, **When** this is detected, **Then** the
   Go binary is NOT produced and the failure is reported clearly.

---

### Edge Cases

- What happens when `pnpm` is not installed but `npm` is? The check MUST report the correct
  required package manager and not silently fall back.
- What happens if the frontend bundle output directory is empty after build? The embedding step
  MUST fail with a diagnostic message rather than silently producing a binary with no UI assets.
- What happens if the `webserver/internal/web/handlers/build/` target directory for embedded
  assets does not exist? The build MUST create it or fail with a clear message.
- What happens on non-Linux platforms? The script MUST either handle or explicitly reject
  unsupported platforms with a clear message.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The build system MUST validate all required tools and their minimum versions
  before executing any build step.
- **FR-002**: The backend MUST be buildable via a single command without any frontend assets
  present.
- **FR-003**: The frontend MUST be buildable via a single command without the backend source
  or toolchain present.
- **FR-004**: The full build MUST produce one self-contained binary with the frontend bundle
  embedded.
- **FR-005**: The build system MUST clean stale artefacts from prior builds before producing
  new ones.
- **FR-006**: Each build stage MUST exit with a non-zero code and a human-readable error on
  failure; no subsequent stage MUST run after a failure.
- **FR-007**: The production binary MUST serve both the frontend UI and backend API from a
  single port when started without the development flag.
- **FR-008**: The build system MUST be invocable from the project root without changing
  directories manually.

### Key Entities

- **Build Stage**: A discrete, ordered step in the pipeline (prerequisite check, frontend
  build, asset copy, backend compile). Each stage MUST be independently runnable and
  independently failable.
- **Embedded Asset Bundle**: The compiled frontend output that is copied into the backend
  source tree and compiled into the binary via the embed directive.
- **Production Binary**: The final output artefact; a single executable containing both the
  API server and the embedded frontend bundle.
- **Prerequisite**: A tool or runtime (Go, Node.js, pnpm) required at a specific minimum
  version for a given build stage.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer on a correctly set up machine can run the full build command and
  obtain a production binary in under 5 minutes.
- **SC-002**: A missing prerequisite is identified and reported within 5 seconds of invoking
  any build command.
- **SC-003**: The backend-only build completes without frontend source or assets present.
- **SC-004**: The frontend-only build completes without the backend toolchain present.
- **SC-005**: The production binary starts and serves the dashboard UI and all API endpoints
  from a single port with no additional files.
- **SC-006**: Build failures at any stage stop the pipeline immediately with a non-zero exit
  code and a message identifying the failed stage.

## Assumptions

- The project runs on Linux; macOS compatibility is considered a nice-to-have but is not in
  scope for this feature.
- `pnpm` is the required frontend package manager; `npm` and `yarn` are not substitutes.
- The minimum Go version is 1.23; the minimum Node.js version is whatever `pnpm` currently
  requires for the project's lockfile.
- The management cluster kubeconfig at `${HOME}/.kube/config` is not required at build time;
  it is only needed at runtime.
- The existing `Makefile` targets are the reference implementation to be reviewed and
  superseded or wrapped by the new build script.
- Container/Docker-based builds are out of scope for this feature; the script targets direct
  host execution.
