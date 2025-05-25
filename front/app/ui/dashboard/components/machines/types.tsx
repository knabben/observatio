import { Meta, Conditions } from "@/app/ui/dashboard/base/types"

/**
 * Represents the details and metadata associated with a machine type.
 */
export type MachineType = {
  metadata: Meta,
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
  metadata: Meta,
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

