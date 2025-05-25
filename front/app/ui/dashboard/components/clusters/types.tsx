import { Meta, Conditions } from "@/app/ui/dashboard/base/types"

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
  metadata: Meta,
  paused: boolean,
  age: string,
  status: {
    ready: boolean,
    phase: string,
    infrastructureReady: boolean,
    controlPlaneReady: boolean,
    conditions: Conditions[]
  }
  clusterNetwork: {
    aPIServerPort: number,
    serviceDomain: string,
    services: {
       cidrBlocks: string[],
       externalIPs: string[],
       nodePortRange: string,
     },
     pods: {
       cidrBlocks: string[],
     },
  }
  controlPlaneEndpoint: {
    host: string,
    port: number,
  }
  controlPlaneRef: {
    apiVersion: string,
    kind: string,
    name: string,
    namespace: string,
  }
  infrastructureRef: {
    apiVersion: string,
    kind: string,
    name: string,
  }
  topology: ClusterClass,
}

export type ClusterInfraType = {
  metadata: Meta,
  cluster: string,
  age: string,
  controlPlaneEndpoint: string,
  server: string,
  thumbprint: string,
  modules: Modules[],
  status: {
    ready: boolean,
    failureReason: string,
    failureMessage: string,
    conditions: Conditions[]
  }
}

export type Modules = {
  controlPlane: boolean,
  targetObjectName: string,
  moduleUUID: string,
}

