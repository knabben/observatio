# Phase 0 Research: Object YAML View & Global AI Troubleshooting Panel

## R1: Does the backend already send the complete object, or only a curated subset?

**Finding**: Every `processor.ProcessX` function (`webserver/internal/infra/clusterapi/processor/*.go`)
embeds the upstream object's full `Status` wholesale (e.g. `cluster.go:43` `Status: cl.Status`), but
reconstructs `Spec` field-by-field into a bespoke DTO, silently dropping many Spec fields. Confirmed
drops include: `ClusterSpec.AvailabilityGates`; most of `VSphereClusterSpec` (CloudProviderConfiguration,
IdentityRef, LoadBalancerRef, ...); `MachineSpec.FailureDomain`/`ReadinessGates`; most of the inlined
VSphere `VirtualMachineCloneSpec` (Snapshot, Datastore, Network, TagIDs, PciDevices, ...);
`MachineDeploymentSpec.Strategy`/`Selector`/`MinReadySeconds`. `ObjectMeta` passes through as the real
`metav1.ObjectMeta` (labels/annotations/finalizers included). WebSocket watchers route through the
same processors — no raw path exists anywhere today (REST or WS).

**Decision**: Add a single new, generic, read-only backend endpoint that bypasses every curated DTO
entirely: `GET /api/raw?group=&version=&resource=&namespace=&name=`, backed by
`clusterapi.NewDynamicClient(ctx).Resource(gvr).Namespace(ns).Get(ctx, name, metav1.GetOptions{})`,
returning `obj.Object` (the full `unstructured.Unstructured` map) as JSON.

**Rationale**: The dynamic client is already used for exactly this style of generic access (Docker
infra fetchers, WS watchers); a `Get` by GVR+namespace+name works identically regardless of Kind, so
one handler serves Cluster, Machine, MachineDeployment, DockerCluster, VSphereCluster, DockerMachine,
and VSphereMachine — no per-type raw handler needed, and existing curated models/handlers are
completely untouched (zero regression risk to Specification tabs or list screens).

**Alternatives considered**:
- Expand every curated model to also embed the full raw `Spec` — rejected: touches shared models used
  by every existing screen and list-stream payload (bandwidth cost for every row in every list, not
  just the one object an operator has open), for a feature only the (rarely-opened) YAML tab needs.
- A per-resource-type raw endpoint (`/api/clusters/raw`, `/api/machines/raw`, ...) — rejected: the
  dynamic-client `Get` is already fully generic by GVR; a single endpoint avoids seven near-identical
  handlers.

## R2: REST vs WebSocket for the raw-object tab (Constitution II)

**Decision**: REST, fetched when the YAML tab is first opened, and re-fetched whenever the
`resourceVersion` of the *already WS-streamed* curated object (feeding the same screen's
Specification tab) changes — not an independent polling loop.

**Rationale**: Building a dedicated WS watcher type for a tab that's opened occasionally, for
inspection, is disproportionate to the value (7 new WS object types for what is fundamentally a
"show me everything" debug view) and Constitution II's live-visibility requirement is already
satisfied by the *existing* per-screen WS stream — this feature makes that same live signal trigger
a REST re-hydration of the one open raw view, rather than replacing the WS delivery mechanism.

**Alternatives considered**: A parallel raw-object WS watcher per Kind (mirroring `WatchDockerClusters`
etc. from feature 004) — rejected as disproportionate; can be revisited later if the tree tab proves
to need true sub-second live updates rather than "updates when the object next changes."

## R3: Tree UI component

**Decision**: Mantine core's built-in `Tree` / `useTree` / `TreeNodeData` — already present in the
pinned `@mantine/core@7.17.8` dependency (confirmed via package inspection: `Tree.d.ts`, `use-tree.d.ts`,
`TreeNode.d.ts` all exist in the installed package; Mantine's own changelog places `Tree`'s
introduction at v7.10.0, before the version already pinned here). No new dependency, no version bump.

**Rationale**: Directly satisfies FR-011 (expandable/collapsible tree) with zero new dependency
surface. A small utility (`to-tree-data.ts`) converts the raw object's arbitrary JSON into
`TreeNodeData[]` (object/array keys become expandable nodes with a `value` path like
`spec.clusterNetwork.pods.cidrBlocks[0]`; scalar leaves render as `"key: value"` labels).

**Alternatives considered**: A dedicated third-party JSON-tree package (e.g. `react-json-view-lite`) —
considered and initially favored before confirming Mantine already ships an equivalent component;
rejected once confirmed, since it would be an unjustified new dependency per the Technology Stack
clause ("No new runtime dependency may be introduced without... justification").

## R4: Collapsible AI panel mechanism

**Decision**: Reuse Mantine's `Drawer` component — already used by `front/app/ui/dashboard/sidenav.tsx`
for the mobile navigation drawer — positioned from the side, mounted once in `dashboard/layout.tsx`
(app-wide), rather than embedded per object type. Open/closed state, the conversation, and the
current auto-context live in a new React context (`ai-panel-context.tsx`), mirroring the
`InfraCapabilityContext` pattern feature `004` already established for exactly this kind of
"fetched/derived once, read by many descendants" state.

**Rationale**: Consistent with an existing, already-tested UI pattern in this codebase (no new
interaction paradigm to design or explain); the context-provider pattern for cross-tree shared state
is likewise already proven in this codebase.

## R5: Auto-context "current object in view" tracking

**Decision**: Each detail screen calls a small hook (`use-current-object-context.ts`) on mount/update
that registers `{kind, name, namespace, status, conditions, keySpecFields}` with the AI panel context,
and unregisters on unmount. The AI panel reads "whatever is currently registered" to build its
pre-fill; a screen with no object in focus (a list view) never registers, so the panel finds nothing
and falls back to an empty/general prompt (FR-007).

**Rationale**: Keeps each screen in control of exactly what it exposes as "key spec fields" (a Cluster
and a Machine have different meaningful fields) without the panel needing per-Kind knowledge baked
into it.

## R6: Conversation persistence scope (Constitution IV)

**Decision**: The AI panel's conversation lives in React state inside the new context provider,
mounted once in the dashboard layout — it survives client-side navigation (the provider doesn't
unmount) but is never written to `localStorage`/`sessionStorage`/any backend store, so a full page
reload clears it, identical to today's embedded-panel behavior.

**Rationale**: Satisfies FR-003 ("closing/reopening within the same session preserves the
conversation") without violating Constitution IV's "MUST NOT persist across sessions without
explicit user consent" — reload = new session, consistent with existing behavior, no new consent UX
needed.
