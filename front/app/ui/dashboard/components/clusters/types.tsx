// Cluster table + details types and models.

type Conditions = {
  type: string,
  status: boolean,
  lastTransitionTime: string,
}

type ClusterClass = {
  isClusterClass: boolean,
  machineDeployments: [],
  kubernetesVersion: string,
}

export type ClusterType = {
  name: string,
  paused: boolean,
  clusterClass: ClusterClass,
  phase: string,
  infrastructureReady: boolean,
  controlPlaneReady: boolean,
  created: string,
  conditions: Conditions[]
}

