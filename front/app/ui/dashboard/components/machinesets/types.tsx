import {Conditions, Meta} from "@/app/ui/dashboard/base/types";

/**
 * Represents the type definition for a MachineSet — replica counts, owning MachineDeployment, and
 * status conditions behind the Day-2 Ops dashboard's stalled-rollout warning (006).
 */
export type MachineSetType = {
  metadata?: Meta,
  age?: string,
  cluster?: string,
  machineDeployment?: string,
  replicas?: number,
  status?: MachineSetStatus,
}

export type MachineSetStatus = {
  selector?: string,
  replicas?: number,
  fullyLabeledReplicas?: number,
  readyReplicas?: number,
  availableReplicas?: number,
  observedGeneration?: number,
  conditions?: Conditions[],
}
