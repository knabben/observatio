// Cluster table + details types and models.

type Conditions = {
  type: string,
  status: boolean,
  lastTransitionTime: string,
}

type MachineDeployments = {
  class: string,
  name: string,
  replicas: string,
  strategy: {
    type: string
  },
}

type ClusterClass = {
  isClusterClass: boolean,
  kubernetesVersion: string,
  className: string,
  controlPlaneReplicas: number,
  machineDeployments: MachineDeployments[]
}

export type ClusterType = {
  name: string,
  paused: boolean,
  podNetwork: string,
  serviceNetwork: string,
  phase: string,
  infrastructureReady: boolean,
  controlPlaneReady: boolean,
  created: string,
  clusterClass: ClusterClass,
  conditions: Conditions[]
}

export type ClusterInfraType = {
  name: string,
  cluster: string,
  created: string,
  server: string,
  ready: boolean
}

