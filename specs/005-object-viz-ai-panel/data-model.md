# Data Model: Object YAML View & Global AI Troubleshooting Panel

No new persistent entities — everything here is either a passthrough of existing cluster data or
transient UI/session state.

## RawObject (backend passthrough, not a Go struct)

`webserver/internal/web/handlers/kubernetes/raw.go`

The new endpoint returns `unstructured.Unstructured.Object` (`map[string]interface{}`) directly —
whatever the Kubernetes API server returns for that GVR/namespace/name, unmodified. No model type is
introduced; there is nothing to keep in sync with upstream CRD changes.

| Query param | Meaning |
|---|---|
| `group` | API group, e.g. `cluster.x-k8s.io`, `infrastructure.cluster.x-k8s.io` |
| `version` | API version, e.g. `v1beta1` |
| `resource` | Plural resource name, e.g. `clusters`, `dockerclusters`, `vspheremachines` |
| `namespace` | Object namespace |
| `name` | Object name |

## TreeNodeData (frontend, Mantine-native type)

`front/app/ui/dashboard/shared/to-tree-data.ts`

Produced by recursively walking the raw object JSON:

| Field | Type | Notes |
|---|---|---|
| `label` | `ReactNode` | `"key"` for expandable object/array nodes, `"key: value"` for scalar leaves |
| `value` | `string` | Unique path, e.g. `spec.clusterNetwork.pods.cidrBlocks[0]` |
| `children` | `TreeNodeData[]?` | Present only for object/array nodes |

## AIPanelState (frontend, in-memory context — not persisted)

`front/app/ui/dashboard/ai-panel/ai-panel-context.tsx`

| Field | Type | Notes |
|---|---|---|
| `isOpen` | `boolean` | Drawer expanded/collapsed |
| `messages` | `WSRequest[]` | Existing conversation shape from `ai-troubleshooting.tsx`, unchanged |
| `currentObjectContext` | `ObjectContext \| null` | Set by whichever detail screen is mounted; `null` on list/overview screens |
| `queryField` | `string` | The editable input; auto-set from `currentObjectContext` only while untouched by the operator |
| `queryFieldTouched` | `boolean` | Set `true` on first manual edit; blocks further auto-refresh until the object changes (FR-008/FR-009) |

## ObjectContext (frontend, derived per screen)

`front/app/ui/dashboard/ai-panel/use-current-object-context.ts`

| Field | Type | Notes |
|---|---|---|
| `kind` | `string` | e.g. `Cluster`, `Machine`, `DockerCluster` |
| `name` | `string` | |
| `namespace` | `string` | |
| `status` | `string` | Human-readable status/condition summary (richer than today's bare condition-reason string) |
| `keySpecFields` | `Record<string, string>` | Per-screen-chosen subset of Spec fields most relevant for troubleshooting that Kind |

## Relationships / Flow

```text
Detail screen mounts
   └─ use-current-object-context() → registers ObjectContext in AIPanelState (global)
   └─ "YAML" tab opened → GET /api/raw?... → to-tree-data() → <Tree data={...}/>
   └─ "Ask AI about this" clicked → AIPanelState.isOpen=true, queryField built from currentObjectContext
Global AIPanelTrigger (any screen) clicked → AIPanelState.isOpen=true, queryField built from
   currentObjectContext if present, else left empty/general
```
