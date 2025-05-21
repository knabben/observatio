/**
 * Represents the details and metadata associated with a machine type.
 */
export type MachineType = {
  metadata: MachineMeta,
  bootstrap: string,
  cluster: string,
  nodeName: string,
  providerID: string,
  version: string,
  age: string,
  status: {
    bootstrapReady: boolean,
    infrastructureReady: boolean,
    conditions: Conditions[]
  }
}


export type MachineInfraType = {
  metadata: MachineMeta,
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

export type MachineMeta = {
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

export type Conditions = {
  type: string,
  reason: string,
  severity: string,
  status: string,
  message: string,
  lastTransitionTime: string,
}
