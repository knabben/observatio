// Machine table + details types and models.

type Conditions = {
  type: string,
  reason: string,
  severity: string,
  status: string,
  message: string,
  lastTransitionTime: string,
}

export type MachineType = {
  name: string,
  namespace: string,
  owner: string,
  bootstrap: string,
  cluster: string,
  nodeName: string,
  providerID: string,
  version: string,
  created: string,
  bootstrapReady: boolean,
  infrastructureReady: boolean,
  phase: string,
}

export type MachineInfraMeta = {
  name: string,
  namespace: string,
  resourceVersion: string,
  creationTimestamp: string,
  labels: {
    [key: string]: string
  },
  annotations: {
    [key: string]: string
  },
  ownerReferences: {
    kind: string,
    name: string,
    uid: string,
    apiVersion: string,
    controller: boolean,
    blockOwnerDeletion: boolean,
  }[]
}

export type MachineInfraType = {
  metadata: MachineInfraMeta,
  providerID: string,
  failureDomain: string,
  powerOffMode: string,
  template: string,
  cloneMode: string,
  numCPUs: number,
  numCoresPerSocket: number,
  memoryMiB: number,
  diskGiB: number,
  age: string,
  status: {
    ready: boolean,
    failureReason: string,
    failureMessage: string,
    conditions: Conditions[]
  }
}
