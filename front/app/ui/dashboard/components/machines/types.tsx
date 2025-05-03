// Cluster table + details types and models.

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
