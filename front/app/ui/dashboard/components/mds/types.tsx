// MachineDeployment models and details

import {Conditions, MachineMeta} from "@/app/ui/dashboard/components/machines/types";

/**
 * Represents the type definition for a MachineDeployment.
 *
 * This type contains metadata and configuration details for a machine deployment,
 * including its replicas, cluster association, age, template versions, and the
 * current deployment status.
 *
 * The properties define specifics such as metadata information, infrastructure
 * references, bootstrap templates, and other essential machine deployment details.
 */
export type MachineDeploymentType = {
  metadata: MachineMeta,
  replicas: number,
  cluster: string,
  age: string,
  templateversion: string,
  templateBootstrap: Bootstrap,
  templateInfrastructureRef: InfrastructureRef,
  status: MachineDeploymentStatus,
}

export type Bootstrap = {
  configRef: {
    name: string,
    namespace: string,
    kind: string,
    apiVersion: string,
  }
}

export type InfrastructureRef = {
  name: string,
  namespace: string,
  kind: string,
  apiVersion: string,
}

export type MachineDeploymentStatus = {
  phase: string,
  conditions: Conditions[],
  readyReplicas: number,
  updatedReplicas: number,
  unavailableReplicas: number,
  v1beta2: {
    availableReplicas: number,
    readyReplicas: number,
    upToDateReplicas: number,
    conditions: Conditions[],
  }
}