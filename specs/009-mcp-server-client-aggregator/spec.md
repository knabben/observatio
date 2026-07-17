# Feature Specification: MCP Server Aggregation & Local Tool Server

**Feature Branch**: `009-mcp-server-client-aggregator`
**Created**: 2026-07-16
**Status**: Draft
**Input**: User description: "make a very strong mcp server and client using the project, it allows agreggating new MCP server to make it even more powerfull, and allows the existent tools to become one of the local mcp servers"

## Clarifications

### Session 2026-07-17

- Q: Should the built-in kubectl-backed capability (FR-001) be packaged as a plain internal
  adapter, or as an actual local MCP server speaking the same protocol as external sources? →
  A: An actual local MCP server, connected in-process (no subprocess) — matching the original
  Input's "become one of the local mcp servers" phrasing, and making it reusable by other MCP
  clients later, not just Observātiō's own aggregator.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Existing tools become a built-in tool source (Priority: P1)

Today, the "Ask AI about this" panel's assistant has exactly one hand-built capability (running
`kubectl` commands). An operator wants that same capability, and any future built-in capability, to be
packaged as a self-contained, standard tool source that the assistant loads the same way it would load
any external tool source — rather than as a one-off, hard-wired special case.

**Why this priority**: Every other story in this feature depends on there being at least one working
tool source to aggregate. Converting the existing capability into that standard shape is the
foundation the aggregation model is built on, and it delivers value on its own: it makes today's one
tool source's behavior (including its safety boundary) explicit and consistent instead of implicit.

**Independent Test**: With no additional tool sources configured, open the "Ask AI about this" panel
on any object and confirm the assistant can still perform the same kubectl-backed inspection it does
today — now visibly presented as one named, self-contained tool source instead of an untitled built-in.

**Acceptance Scenarios**:

1. **Given** no additional tool sources have been configured, **When** an operator opens the "Ask AI
   about this" panel, **Then** the assistant offers the same inspection capability available today,
   identified as a distinct, named tool source.
2. **Given** the built-in tool source is temporarily unavailable, **When** an operator asks the
   assistant a question, **Then** the assistant reports that the capability is unavailable rather than
   silently failing or hanging.

---

### User Story 2 - Add a new tool source without changing code (Priority: P1)

An operator wants to give the assistant new abilities — for example, real backup/restore verbs from a
community-maintained tool source — by registering that tool source, not by waiting for a code change
and a new release of the product.

**Why this priority**: This is the "aggregating new MCP server" half of the request and the main
promise of the feature: the assistant's capability set grows by configuration, not by development.
Without it, the feature is just a rename of the existing single tool.

**Independent Test**: Register one additional external tool source pointing at a running, reachable
tool server; without any other change, confirm the assistant's next answer can draw on a capability
that source provides and that wasn't available before the registration.

**Acceptance Scenarios**:

1. **Given** a reachable external tool source is registered, **When** an operator asks the assistant a
   question that capability can answer, **Then** the assistant uses it and produces a response backed
   by that tool's output.
2. **Given** two tool sources are registered and both are reachable, **When** the assistant is asked a
   question, **Then** it can draw on capabilities from either source in the same conversation.
3. **Given** an operator removes or disables a previously-registered tool source, **When** the
   assistant is next asked a question, **Then** it no longer offers or uses that source's capabilities.
4. **Given** a newly registered tool source exposes a capability with the same name as one already
   offered by another registered source, **When** registration completes, **Then** the conflict is
   surfaced to the operator rather than silently resolved, and the assistant's tool list stays
   unambiguous.

---

### User Story 3 - Aggregate resilience when a source misbehaves or disappears (Priority: P2)

An operator has several tool sources registered (the built-in one plus one or more external ones). One
of the external sources becomes slow, unreachable, or starts erroring. The operator expects the
assistant to keep working with whatever sources are still healthy, and to be able to see which source
is degraded.

**Why this priority**: Aggregating multiple independent sources only becomes trustworthy in practice
once a bad source can't take the whole assistant down with it. This is what makes the aggregation
"powerful" rather than "fragile," but the product is still useful with a single source if this isn't
built yet.

**Independent Test**: With two tool sources registered, make one unreachable (e.g., stop it) and
confirm the assistant still answers questions using the remaining healthy source, while the unhealthy
one is visibly flagged somewhere an operator would look.

**Acceptance Scenarios**:

1. **Given** one of several registered tool sources stops responding, **When** the assistant is asked
   a question answerable by a healthy source, **Then** it answers normally without being blocked by
   the unhealthy one.
2. **Given** a registered tool source is unreachable, **When** an operator checks the status of
   registered tool sources, **Then** that source is shown as unhealthy, with the others shown as
   healthy.
3. **Given** a previously-unhealthy tool source becomes reachable again, **When** the assistant is next
   asked a question, **Then** its capabilities are available again without requiring the operator to
   re-register it.

---

### Edge Cases

- What happens when two registered tool sources offer conflicting or contradictory answers for
  overlapping capabilities (e.g., both claim to describe backup health)? The assistant's response
  should make clear which source it drew the answer from.
- What happens when an operator registers a tool source that turns out to expose a mutating/write
  capability (e.g., "trigger a restore")? See the read-only boundary requirement below.
- What happens when a tool source is reachable but returns malformed or unexpected data? The assistant
  should surface a clear error for that capability rather than presenting corrupted output as fact.
- What happens when every registered tool source is unreachable at once? The assistant should report
  that it currently has no working capabilities rather than appearing to hang.
- What happens when the same underlying action (e.g., "list backups") is offered by both the built-in
  source and an external one? The operator needs a way to tell them apart when reviewing registered
  sources.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The assistant MUST offer its existing kubectl-backed inspection capability as a
  self-contained, named tool source, functionally equivalent to today's behavior, rather than as a
  hard-wired special case. This tool source MUST speak the same protocol as every external tool
  source (Clarifications, 2026-07-17) — an internal adapter that merely mimics one is not
  sufficient — so the built-in capability is structurally indistinguishable from a registered
  external one anywhere the aggregator or an operator inspects it.
- **FR-002**: Operators MUST be able to register an additional external tool source by supplying its
  connection information, without requiring a code change or new release.
- **FR-003**: Operators MUST be able to view the list of currently registered tool sources, each
  source's health/reachability, and which capabilities it currently contributes.
- **FR-004**: Operators MUST be able to disable or remove a previously-registered tool source, after
  which the assistant no longer offers or uses its capabilities.
- **FR-005**: The assistant MUST present the union of capabilities from all currently healthy,
  registered tool sources as a single, unified capability set within one conversation — an operator
  should not need to know which source backs which capability to use it.
- **FR-006**: When a registered tool source is unreachable or erroring, the assistant MUST continue
  operating using the remaining healthy sources rather than failing the whole conversation.
- **FR-007**: When two registered tool sources offer a capability with the same name, the system MUST
  surface that naming conflict to the operator at registration time rather than silently picking one or
  merging them.
- **FR-008**: Every response the assistant gives MUST make it possible to identify which tool source(s)
  contributed the underlying data, for traceability and trust.
- **FR-009**: All capabilities offered through any registered tool source — built-in or external — MUST
  be treated as read-only/non-mutating from the assistant's perspective; the assistant MUST NOT execute
  a capability that would change cluster or backup state.
- **FR-010**: Registering, modifying, or removing a tool source MUST be restricted to an operator
  performing administrative configuration of Observātiō, not an end-user action available from within
  the "Ask AI about this" chat panel itself.

### Key Entities

- **Tool Source**: A registered provider of assistant capabilities — either the built-in one wrapping
  today's existing capability, or an external one an operator has added. Has a name, connection
  information, a health/reachability status, and the set of capabilities it currently contributes.
- **Capability**: A single named action the assistant can invoke (e.g., "run a read-only kubectl
  command," "describe a backup"), contributed by exactly one tool source at a time and identified by
  which source it came from.
- **Assistant Conversation**: An "Ask AI about this" session in which the assistant draws on the
  aggregated set of capabilities from all healthy, registered tool sources to answer a question about
  an object.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator can add a new capability to the assistant (by registering an external tool
  source) in under 5 minutes, with no code change or redeployment of Observātiō itself.
- **SC-002**: With any single registered tool source unreachable, the assistant still successfully
  answers 100% of questions that depend only on the remaining healthy sources.
- **SC-003**: An operator can always determine, for any assistant answer, which registered tool
  source(s) it came from.
- **SC-004**: 100% of capabilities exposed through the assistant, across all registered tool sources,
  are verified non-mutating before being made callable.
- **SC-005**: An operator can view the health of all registered tool sources and identify an unhealthy
  one within a few seconds, without inspecting logs.

## Assumptions

- The built-in kubectl-backed capability, once wrapped as the "local" tool source, keeps its current
  behavior and constraints; this feature does not change what that specific capability can do, only how
  it's packaged and combined with others. Per Clarifications (2026-07-17), it is packaged as an actual
  local MCP server rather than a plain internal adapter — a side effect of that choice is that this
  same local server could, in principle, be connected to by other MCP clients later, though doing so is
  not built or exposed by this feature (see the out-of-scope bullet below, which still stands).
- This feature scopes only the assistant's outbound tool-consumption side (Observātiō as an aggregator/
  client of tool sources) for the existing "Ask AI about this" panel. Exposing Observātiō's own Day-2
  Ops/backup-health state as a tool source *for other, external assistants to consume* (the reverse
  direction raised in prior proposals) is a distinct, larger effort and is out of scope here.
- Realistic external tool sources this feature targets include in-cluster MCP servers reachable over
  HTTP (e.g. `velero-mcp`, which runs as a workload inside the management cluster and exposes
  backup/restore capabilities) and MCP servers built or hosted via toolkits such as `kmcp` — both are
  ordinary registered external sources under FR-002, requiring no source-specific handling.
- Tool source registration is an administrative/deployment-time concern (config-driven), not a
  per-conversation or per-user runtime action, since this mirrors how the rest of Observātiō is
  configured today and keeps the safety boundary under operator control.
- Enforcing the read-only boundary applies to whether a capability may be *invoked* through the
  assistant; it does not require modifying or restricting what a registered external tool source itself
  is capable of outside this integration.
- A reasonable number of simultaneously registered tool sources (e.g., low tens) is assumed; the feature
  is not required to scale to hundreds of concurrently aggregated sources for this iteration.
