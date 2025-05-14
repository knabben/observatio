// Machine table + details types and models.

type Conditions = {
  type: string,
  status: string,
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

export type MachineInfraType = {
  name: string,
  namespace: string,
  providerID: string,
  failureDomain: string,
  powerOffMode: string,
  template: string,
  cloneMode: string,
  numCPUs: number,
  numCoresPerSocket: number,
  memoryMiB: number,
  diskGiB: number,
  ready: boolean,
  failureReason: string,
  failureMessage: string,
  created: string,
  conditions: Conditions[]
}
