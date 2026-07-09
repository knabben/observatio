# Contract: New `ObjectType`s and GVRs

Extends the existing `/ws/watcher` dispatch table (`webserver/internal/web/handlers/system/
websocket.go`'s `watchHandlers` map) and `front/app/lib/resource-gvr.ts`'s `RESOURCE_GVR` with four
new entries, following exactly the shape every existing kind already uses (no new envelope, no new
transport — see `webserver/internal/web/watchers/processor.go`'s `EventResponse{Type, Event, Data}`).

| `ObjectType` (WS `{"type": ...}`) | GVR | Backend watcher | Frontend route |
|---|---|---|---|
| `machinehealthcheck` | `cluster.x-k8s.io/v1beta1, Resource=machinehealthchecks` | `watchers.WatchMachineHealthChecks` | `/dashboard/machinehealthchecks` |
| `kubeadmcontrolplane` | `controlplane.cluster.x-k8s.io/v1beta1, Resource=kubeadmcontrolplanes` | `watchers.WatchKubeadmControlPlanes` | `/dashboard/kubeadmcontrolplanes` |
| `machineset` | `cluster.x-k8s.io/v1beta1, Resource=machinesets` | `watchers.WatchMachineSets` | `/dashboard/machinesets` |
| `clusterclass` | `cluster.x-k8s.io/v1beta1, Resource=clusterclasses` | `watchers.WatchClusterClasses` | `/dashboard/clusterclasses` |

Each watcher follows the identical 1:1 relay shape as every existing one (`machine.go`'s
`WatchMachines`): open one `dynamicClient.Resource(gvr).Namespace("").Watch(...)`, convert each
event via the kind's `processor.Process<Kind>`, and stream it through
`watchers.WatchResourceViaWebSocket`. No fan-in, no aggregation — unlike 006's `day2ops.go`, these
are single-purpose list-page watchers.

## YAML tab (unchanged endpoint, new query values)

`GET /api/raw?group=&version=&resource=&namespace=&name=` (established in 005) is reused as-is; the
four new kinds' detail screens simply pass their `RESOURCE_GVR` entry — no backend change required.

## Empty/unavailable CRD handling (FR-012)

When a watch fails to open because the target CRD isn't installed (e.g. `kubeadmcontrolplanes` on a
non-kubeadm control-plane provider), the watcher returns an error the same way every existing
watcher already does on a genuine failure — `HandleWatcher` logs it and the connection closes
without ever sending data. After the frontend's bounded reconnect attempts are exhausted,
`useResourceStream` settles into its existing `error` `ChannelState`, which `BaseLister`
(`front/app/ui/dashboard/base/lister.tsx`) already renders as "Unable to load `<kind>`. The
connection may be unavailable." with a Retry action — not a raw/unhandled error or a blank screen.
This already satisfies FR-012 with no new frontend logic; when the kind's objects simply don't
exist yet (CRD present, zero items), the same hook instead settles into its `empty` state, rendered
as "No `<kind>` found."
