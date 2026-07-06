# Feature Specification: Object YAML View & Global AI Troubleshooting Panel

**Feature Branch**: `005-object-viz-ai-panel`
**Created**: 2026-07-06
**Status**: Draft
**Input**: User description: "we need IMPROVEMENTS in the object visualization details for all component clusters, mds, etc. This is going to be achieved adding a new panel to print the entire yaml object as a tree in another tab on each details screen, this allows a another visualization of the actual object. the second improvement is related to the AI panel, and existent stream, the context is very poor and need to be increased, the screen colors are not matching also, the panel is not an entire section instead must be a colapsable screen that can be called anytime from the side. and the context of the current screen object can be added automatically in the field before searching."

## Overview

Two related but independent gaps in how operators inspect objects and get AI help across
Observātiō's detail screens (Clusters, Cluster Infrastructure, Machines, Machine Infrastructure,
Machine Deployments):

1. **Raw object visualization**: every detail screen renders only the specific fields a given
   component author chose to surface. Operators who need to see the complete underlying object —
   to check a field nobody built a UI for, or to compare against what they applied — have no way to
   do so today; they'd have to leave the dashboard and query the cluster directly.
2. **AI troubleshooting is bolted onto one tab, per object type, with poor context and mismatched
   styling**: it's embedded as a full-width section inside a specific object's detail tabs
   (available only when that object is open), pre-fills its query with little more than a
   condition-reason string, and uses a hardcoded dark color scheme that doesn't match the rest of
   the dashboard's theme. Operators can't bring the assistant with them as they navigate, and when
   they do use it, it starts from thin context.

This feature makes the complete object available as a readable tree on every detail screen, and
turns AI troubleshooting into a persistent, collapsible panel reachable from anywhere in the
dashboard, pre-loaded with rich context about whatever the operator is currently looking at, and
styled consistently with the rest of the app.

## Clarifications

### Session 2026-07-06

- Q: With the new full-object tree tab added, should the tab layout stay at 2 tabs (Specification +
  Tree), or should curated tables be reduced/removed since the tree already shows everything? → A:
  Keep 2 tabs — Specification (curated summary + conditions table, unchanged) and the new Tree tab.
  No separate "Status" tab is introduced, since that would duplicate what's already in
  Specification's header/conditions table. Object status MUST remain prominent within Specification
  (header status indicator + conditions table) — consolidating tabs must not bury or diminish it.
- Q: Should each object's detail screen get a one-click "Ask AI about this" action that jumps
  straight into the global panel pre-loaded with that object, in addition to the passive global
  side-trigger? → A: Yes — add a per-object quick-action so opening AI help on the object currently
  in view is a single click, not "open the global panel and notice it auto-filled."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - AI troubleshooting is available from anywhere, not just one tab (Priority: P1)

An operator is looking at a Machine's detail screen (or any other screen in the dashboard) and
wants AI help. Today they'd only have that option if they're on a specific object's detail view
with the "AI Troubleshooting" tab selected. Instead, they can open a collapsible AI panel from a
persistent control at any time, from any screen, without navigating away from what they're doing.

**Why this priority**: This is the core structural fix — an assistant that's only reachable from
one tab of one screen type isn't actually available "anytime," which undermines the whole feature.
Every other improvement to the AI panel depends on it existing as a global, reachable surface
first.

**Independent Test**: From the Clusters list, the Dashboard overview, and a Machine detail screen,
confirm a control to open the AI panel is visible and opens the same collapsible panel each time;
confirm the panel can be collapsed and reopened without losing anything already typed or received
in the current session.

**Acceptance Scenarios**:

1. **Given** any dashboard screen, **When** the operator activates the AI panel control, **Then** a
   collapsible panel slides in from the side without navigating away from the current screen.
2. **Given** the AI panel is open, **When** the operator activates the collapse control, **Then**
   the panel closes and the underlying screen is fully usable again.
3. **Given** the AI panel was open with an in-progress conversation, **When** the operator
   collapses and later reopens it (without a full page reload), **Then** the prior conversation is
   still visible.
4. **Given** an object's detail screen, **When** the operator opens it, **Then** no "AI
   Troubleshooting" tab appears among that object's detail tabs — the global panel is the only
   entry point to AI troubleshooting.

---

### User Story 2 - The AI panel starts from rich, automatic context (Priority: P1)

An operator opens the AI panel while viewing a specific object (e.g., a failing Machine). Today the
pre-filled query is little more than a condition-reason string. Instead, the query field is
automatically populated with a substantially fuller description of the object currently in view —
enough that the operator can send it as-is and get a useful answer, or edit it before sending.

**Why this priority**: Poor context is called out explicitly as a core problem; a panel that's
reachable everywhere but still gives the assistant almost nothing to work with only half-solves the
problem. This depends on User Story 1's global panel existing, but is independently verifiable once
it does.

**Acceptance Scenarios**:

1. **Given** the operator is viewing a specific object's detail screen, **When** they open the AI
   panel, **Then** the query field is pre-filled with a description covering that object's
   identity, current status/conditions, and key spec fields — not just a bare condition reason.
2. **Given** the operator is on a screen with no single object in focus (e.g., a list view or the
   Dashboard overview), **When** they open the AI panel, **Then** the query field is empty or
   contains a general prompt rather than a broken/partial context string.
3. **Given** the pre-filled context, **When** the operator edits or replaces it before sending,
   **Then** their edit is what gets sent, not the auto-filled text.
4. **Given** the operator navigates to a different object's detail screen while the AI panel is
   already open, **When** the panel's query field is still empty (not yet edited or sent), **Then**
   it updates to reflect the newly-viewed object's context.
5. **Given** an object's detail screen, **When** the operator activates that screen's "Ask AI about
   this" quick-action, **Then** the global AI panel opens (if not already open) pre-filled with that
   object's context in one click, without first requiring the operator to open the panel separately.

---

### User Story 3 - Inspect the complete raw object on any detail screen (Priority: P2)

An operator viewing any object's detail screen (Cluster, Cluster Infrastructure, Machine, Machine
Infrastructure, Machine Deployment) wants to see the complete underlying object, not just the
fields the screen chose to surface. They open a new tab on that same detail screen and see the
entire object rendered as a readable, navigable tree.

**Why this priority**: This is a valuable, additive inspection capability, but — unlike the AI
panel fixes — it doesn't correct an existing broken/inconsistent experience, so it's ranked after
the AI panel corrections.

**Independent Test**: Open the detail view for a Cluster, a Machine, and a Machine Deployment; on
each, open the new tab and confirm the complete object (every field the backend returns for it, not
a curated subset) is visible in a readable, expandable/collapsible tree structure.

**Acceptance Scenarios**:

1. **Given** any object's detail screen, **When** the operator selects the new visualization tab,
   **Then** the complete object is rendered as a tree with expandable/collapsible nodes for nested
   fields.
2. **Given** the rendered tree, **When** the operator collapses a nested section (e.g., `status`),
   **Then** that section's children are hidden and can be re-expanded.
3. **Given** an object with a very large or deeply nested field (e.g., long condition history),
   **When** it is rendered, **Then** the tab remains responsive and scrollable rather than freezing
   or overflowing the page.
4. **Given** an object whose data updates while the tab is open (live stream), **When** new data
   arrives, **Then** the tree reflects the updated object rather than showing stale data
   indefinitely.

---

### Edge Cases

- What happens if the AI panel is opened while the underlying live stream for the current object
  has failed or is disconnected? (The panel should still open; context pre-fill uses the last known
  data, or is omitted if none is available yet — it never crashes the panel.)
- What happens if an object has no meaningful conditions or status yet (e.g., just created)? (The
  auto-filled context states that plainly rather than fabricating detail.)
- What happens on very small viewports (narrow browser windows) when the AI panel is expanded?
  (The panel must remain usable and not permanently obscure the entire screen without a way back.)
- What happens if the operator sends a message, then navigates to a different screen before a
  response arrives? (The response still arrives and appears in the panel/conversation when it does,
  since the panel persists across navigation.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a persistent control, reachable from any dashboard screen,
  that opens a collapsible AI troubleshooting panel from the side.
- **FR-002**: The AI panel MUST be collapsible/expandable without navigating away from or losing
  the state of the underlying screen.
- **FR-003**: Closing and reopening the AI panel within the same session MUST preserve the current
  conversation (messages already sent/received are not lost).
- **FR-004**: The embedded "AI Troubleshooting" tab currently present on each object's detail
  screen MUST be removed; the global panel becomes the sole entry point to AI troubleshooting.
- **FR-005**: The "Object conditions" table currently shown alongside the embedded AI panel MUST
  remain available on the object's detail screen (folded into the "Specification" tab) after the
  AI Troubleshooting tab is removed. No separate "Status" tab is introduced — object status
  (header indicator + conditions table) MUST stay prominent within the Specification tab rather
  than being diminished or buried by consolidation.
- **FR-006**: When the AI panel is opened while a specific object's detail screen is in view, the
  system MUST automatically populate the query field with a description covering that object's
  identity (name/namespace/type), current status/conditions, and key specification fields.
- **FR-007**: When the AI panel is opened from a screen with no single object in focus, the system
  MUST leave the query field empty or show a general prompt rather than a broken partial-context
  string.
- **FR-008**: The operator MUST be able to freely edit or fully replace the automatically populated
  context before sending it; whatever is in the field at send time is what gets sent.
- **FR-009**: If the currently-viewed object changes while the AI panel is open and its query field
  has not yet been edited or sent, the system MUST refresh the pre-filled context to match the
  newly-viewed object.
- **FR-010**: The AI panel's visual styling (colors, contrast) MUST use the dashboard's existing
  theme tokens instead of the current hardcoded, mismatched color values.
- **FR-011**: Every object detail screen (Cluster, Cluster Infrastructure — both Docker and
  vSphere, Machine, Machine Infrastructure — both Docker and vSphere, Machine Deployment) MUST gain
  a new tab that renders the complete underlying object as an expandable/collapsible tree.
- **FR-012**: The rendered object tree MUST reflect every field the backend returns for that
  object, not a curated subset chosen by the existing screen-specific components.
- **FR-013**: The object tree view MUST remain responsive and scrollable for large/deeply nested
  objects rather than freezing or overflowing the page.
- **FR-014**: When the underlying object's live data updates while the tree tab is open, the tree
  MUST reflect the update rather than remaining stale indefinitely.
- **FR-015**: The AI panel MUST remain usable (not permanently obscuring the screen with no way to
  close it) at narrow/mobile viewport widths.
- **FR-016**: Each object's detail screen MUST offer a one-click "Ask AI about this" quick-action
  that opens the global AI panel already pre-filled with that object's context, in addition to the
  passive global side-trigger (FR-001) that pre-fills based on whatever is currently in view.

### Key Entities

- **AI Troubleshooting Panel**: A single, app-wide collapsible UI surface (not tied to any one
  object type) holding the current conversation and its expand/collapse state.
- **Auto-Populated Context**: A derived, editable text description of "whatever object is currently
  in view," built from that object's identity, status/conditions, and key spec fields — refreshed
  when the viewed object changes and the field hasn't been touched.
- **Object Tree View**: A read-only, expandable/collapsible rendering of an object's complete data,
  available as a tab alongside each object's existing detail tabs.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operators can open the AI panel from any dashboard screen in one action, without
  navigating to a specific object first.
- **SC-002**: 100% of object detail screens (Cluster, Cluster Infrastructure, Machine, Machine
  Infrastructure, Machine Deployment) offer a complete-object tree view tab.
- **SC-003**: When opened from an object's detail screen, the AI panel's pre-filled context
  includes the object's identity, status, and key spec fields, not only a condition-reason phrase.
- **SC-004**: The AI panel's visual styling passes the same contrast/theme-consistency bar as the
  rest of the dashboard (no hardcoded colors that clash with the active theme).
- **SC-005**: Collapsing and reopening the AI panel never loses the operator's current conversation
  within a session.
- **SC-006**: The object tree view renders without freezing or breaking layout for the largest
  objects currently seen in production (deeply nested status/condition history).
- **SC-007**: From any object's detail screen, operators can reach a pre-filled AI panel in exactly
  one click via that screen's quick-action, without a separate step to open the panel first.

## Assumptions

- "Print the entire yaml object as a tree" is interpreted as a structured, expandable/collapsible
  tree rendering of the object's data (equivalent in content to its YAML representation), not a
  literal YAML-syntax text block — this gives operators a more navigable view of deeply nested
  fields than a flat text dump.
- "Richer context" means expanding the AI query pre-fill to the object's identity, status,
  conditions, and key spec fields — it does not include fetching controller logs; a deeper
  cluster-detail/log-viewing capability was already scoped out as a separate, dedicated follow-on
  feature in a prior session and remains out of scope here.
- "Called anytime from the side" is interpreted as a persistent trigger available in the dashboard's
  global UI shell (e.g., navigation/side area), reachable regardless of which screen is active; each
  object detail screen additionally gets its own one-click "Ask AI about this" quick-action
  (FR-016) for tighter, more direct integration than the passive global trigger alone.
- The AI panel's conversation is a single, app-wide thread (not one conversation per object); moving
  between objects changes what auto-context would be offered next, but does not itself clear
  history already in the conversation.
- Object tree data comes from data already available to the frontend for the currently-viewed
  object (the same live-streamed object backing the existing detail fields) — this feature does not
  require any new backend endpoint to fetch additional object data.
