# Feature Specification: Infrastructure Provider Detection & Adaptive Listing Screens

**Feature Branch**: `004-detect-infra-adapt-ui`
**Created**: 2026-07-05
**Status**: Draft
**Input**: User description: "the project must be able to detect the infrastructure cluster (docker/vsphere) supported and adapt the listing screens for this"

## Overview

Observātiō currently assumes every connected Cluster API environment is backed by vSphere: the Clusters screen always shows a static "vSphere Clusters" tab and vSphere-specific fields (server, thumbprint, modules), regardless of what infrastructure actually backs the clusters in view. Environments backed by Docker (common for local/dev/test Cluster API setups) have no equivalent view today, and would otherwise show a vSphere tab with no data or a misleading empty state.

This feature makes the dashboard detect, from the connected cluster's Cluster API resources, which infrastructure provider(s) — Docker, vSphere, or both — actually back the clusters it is displaying, and adapts the listing screens accordingly: showing the right provider-specific view (and hiding/relabeling the wrong one), and giving operators a clear, per-cluster indicator of which provider is in play. The connected cluster context is assumed to already carry the Cluster API CRDs needed to identify and inspect each cluster's infrastructure, so detection is derived from data already available rather than requiring new configuration.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Listing screens reflect the actual infrastructure provider in use (Priority: P1)

An operator connects the dashboard to a Cluster API environment backed by Docker. Today they would see a "vSphere Clusters" tab that is empty or irrelevant. Instead, the dashboard should recognize that the clusters present are Docker-backed and present a Docker-appropriate infrastructure view, with no vSphere-only tab confusing the picture.

**Why this priority**: This is the core problem statement — showing the wrong (or no) provider view actively misleads operators about what infrastructure they're looking at, and blocks any real usage of the dashboard against non-vSphere environments.

**Independent Test**: Point the dashboard at a Docker-backed Cluster API environment and confirm the Clusters screen shows a Docker-appropriate infrastructure view (not a vSphere-labeled empty tab); point it at a vSphere-backed environment and confirm the existing vSphere view still appears correctly.

**Acceptance Scenarios**:

1. **Given** a connected environment where all clusters are Docker-backed, **When** the operator opens the Clusters screen, **Then** the infrastructure detail view shown is Docker-specific and no vSphere-only tab/screen is presented as if applicable.
2. **Given** a connected environment where all clusters are vSphere-backed, **When** the operator opens the Clusters screen, **Then** the existing vSphere infrastructure view (server, thumbprint, modules) continues to render exactly as before.
3. **Given** a connected environment with both Docker-backed and vSphere-backed clusters present, **When** the operator opens the Clusters screen, **Then** both provider-specific views are available and each shows only the clusters that belong to it.

---

### User Story 2 - Operators can identify a cluster's infrastructure provider at a glance (Priority: P2)

While scanning the main Clusters list, an operator wants to know immediately whether a given cluster is Docker- or vSphere-backed, without opening a detail view or a separate tab.

**Why this priority**: Faster comprehension of the fleet reduces the chance an operator applies the wrong mental model (e.g., expecting vSphere fields on a Docker cluster) and speeds up troubleshooting across mixed environments.

**Independent Test**: Load the main Clusters list against a mixed environment and confirm every row shows a clear, correct provider indicator matching that cluster's actual backing infrastructure.

**Acceptance Scenarios**:

1. **Given** the main Clusters list with both Docker- and vSphere-backed clusters, **When** it renders, **Then** each row displays an indicator identifying its infrastructure provider.
2. **Given** a cluster whose infrastructure provider cannot be determined, **When** it appears in the list, **Then** it shows an "unknown/unsupported provider" indicator rather than a blank or incorrect one.

---

### User Story 3 - Unsupported or undetectable providers degrade gracefully (Priority: P3)

An operator connects the dashboard to an environment using a Cluster API infrastructure provider other than Docker or vSphere, or to a cluster resource whose infrastructure reference cannot be resolved.

**Why this priority**: The dashboard must not crash, hang, or silently drop clusters just because a provider outside the currently supported set is present — graceful degradation preserves trust and keeps the rest of the dashboard usable.

**Independent Test**: Point the dashboard at (or inject) a cluster resource with an infrastructure reference to a provider other than Docker/vSphere and confirm it still appears in the main list with a generic/unknown indicator and no provider-specific view is forced onto it.

**Acceptance Scenarios**:

1. **Given** a cluster backed by an infrastructure provider outside the supported set, **When** the Clusters screen loads, **Then** the cluster still appears in the main list with a generic/unknown provider indicator and no crash occurs.
2. **Given** an environment where no supported infrastructure provider CRDs are detected at all, **When** the operator opens the Clusters screen, **Then** a clear message communicates that no supported infrastructure provider was detected, instead of an empty screen with no explanation.

---

### Edge Cases

- What happens when a management cluster has both Docker and vSphere provider CRDs installed, but zero clusters currently exist for one of them? (The corresponding view should show a correct "no clusters" empty state, not be hidden entirely, since the provider is genuinely supported here.)
- How does the system handle a cluster resource whose infrastructure reference is missing or points to a resource that no longer exists?
- How does the system handle a Cluster API environment where the connected context lacks the CRDs to identify infrastructure at all (bare core Cluster API only, no infrastructure provider installed)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST determine, for each cluster resource in the connected environment, which infrastructure provider (Docker or vSphere) backs it, based on that cluster's infrastructure reference.
- **FR-002**: System MUST determine, at the environment level, whether Docker infrastructure support, vSphere infrastructure support, or both are present, before deciding which provider-specific listing screens/tabs to offer.
- **FR-003**: The Clusters listing screen MUST display, per cluster row, an indicator of which infrastructure provider backs that cluster.
- **FR-004**: A provider-specific infrastructure listing view MUST only be shown as applicable when that provider is actually detected in the connected environment; it MUST NOT be hardcoded to always appear regardless of what's actually present.
- **FR-005**: When Docker infrastructure is detected, the system MUST provide a Docker-appropriate infrastructure listing view showing the fields relevant to Docker-backed clusters, equivalent in role to the existing vSphere infrastructure view.
- **FR-006**: When a cluster's infrastructure provider is neither Docker nor vSphere, the system MUST still show that cluster in the main list with a generic/unknown-provider indicator, rather than omitting it or crashing.
- **FR-007**: Provider detection MUST happen automatically from data already available in the connected cluster; operators MUST NOT need to manually select or configure which provider is in use.
- **FR-008**: The Machines and Machine Deployments listing screens MUST reflect the same detected-provider adaptation as the Clusters screen wherever they currently expose provider-specific fields.
- **FR-009**: If no supported infrastructure provider is detected in the connected environment, the system MUST present a clear message stating that, instead of leaving a screen empty with no explanation.
- **FR-010**: Existing vSphere-specific listing behavior and fields MUST continue to work unchanged for vSphere-backed clusters; this feature changes when/how that view is surfaced, not what it shows.

### Key Entities

- **Cluster**: Existing entity representing a Cluster API cluster; gains a derived "Infrastructure Provider" attribute (Docker, vSphere, or Unknown/Unsupported) resolved from its infrastructure reference.
- **Infrastructure Provider Capability**: A per-environment detection result indicating whether Docker support, vSphere support, both, or neither are present among the connected cluster's Cluster API resources.
- **Cluster Infrastructure Details**: Provider-specific detail data associated with a cluster (e.g., vSphere: server, thumbprint, modules; Docker: the equivalent set of Docker-relevant infrastructure attributes).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: When connected to a Docker-only environment, operators no longer see a vSphere-labeled tab or screen with no relevant data.
- **SC-002**: 100% of clusters shown in the Clusters list display a provider indicator (Docker, vSphere, or Unknown) that matches that cluster's actual backing infrastructure.
- **SC-003**: Operators can identify which infrastructure provider backs any given cluster within one glance at the listing screen, without needing to open a detail view.
- **SC-004**: The dashboard renders a fully working, non-crashing Clusters screen when evaluated against a Docker-only environment, a vSphere-only environment, and a mixed environment.
- **SC-005**: Detecting the infrastructure provider requires zero manual configuration or selection steps from the operator.

## Assumptions

- The Kubernetes context the dashboard connects to already has the Cluster API CRDs installed for whichever infrastructure provider(s) are in use, sufficient to list clusters and inspect their infrastructure details — detection is derived from these existing resources rather than any new external configuration.
- Docker and vSphere are the two infrastructure providers in scope for adaptive views in this feature; other Cluster API infrastructure providers (e.g., AWS, Azure) are out of scope for building a dedicated view, but must still degrade gracefully to a generic/unknown indicator rather than causing errors.
- A given cluster resource is backed by exactly one infrastructure provider at a time (per its infrastructure reference), though the connected environment may have multiple providers' CRDs installed simultaneously (a mixed management cluster).
- Existing vSphere-specific functionality remains fully available for vSphere-backed clusters; this feature changes how and when that view is surfaced, not its underlying content.
