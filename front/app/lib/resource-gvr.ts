import {ResourceGVR} from '@/app/lib/data';

/**
 * GroupVersionResource for each object kind the dashboard renders a detail screen for, used to
 * request the complete raw object via GET /api/raw (contracts/raw-object-api.md). Mirrors the
 * GVR constants already defined server-side in webserver/internal/web/watchers/*.go.
 */
export const RESOURCE_GVR = {
  cluster: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'clusters'},
  vsphereCluster: {group: 'infrastructure.cluster.x-k8s.io', version: 'v1beta1', resource: 'vsphereclusters'},
  dockerCluster: {group: 'infrastructure.cluster.x-k8s.io', version: 'v1beta1', resource: 'dockerclusters'},
  machine: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'machines'},
  vsphereMachine: {group: 'infrastructure.cluster.x-k8s.io', version: 'v1beta1', resource: 'vspheremachines'},
  dockerMachine: {group: 'infrastructure.cluster.x-k8s.io', version: 'v1beta1', resource: 'dockermachines'},
  machineDeployment: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'machinedeployments'},
} as const satisfies Record<string, ResourceGVR>;
