// Cluster table + details types and models.

type Conditions = {
  type: string,
  status: string,
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
  namespace: string,
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

export type Modules = {
  controlPlane: boolean,
  targetObjectName: string,
  moduleUUID: string,
}

export type ClusterInfraType = {
  name: string,
  namespace: string,
  cluster: string,
  created: string,
  controlPlaneEndpoint: string,
  server: string,
  thumbprint: string,
  ready: boolean,
  modules: Modules[],
  conditions: Conditions[]
}

