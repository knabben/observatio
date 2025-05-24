// MachineDeployment models and details

import {Conditions, MachineMeta} from "@/app/ui/dashboard/components/machines/types";

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
  condition: Conditions[],
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