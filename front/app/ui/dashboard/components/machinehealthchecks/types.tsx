import {Conditions, Meta} from "@/app/ui/dashboard/base/types";

/**
 * Represents the type definition for a MachineHealthCheck — the remediation policy (target
 * selector, timeouts, maxUnhealthy threshold) behind the Day-2 Ops dashboard's self-healing/
 * needs-investigation severity classification (006/US4).
 */
export type MachineHealthCheckType = {
  metadata?: Meta,
  age?: string,
  cluster?: string,
  selector?: Selector,
  maxUnhealthy?: string,
  nodeStartupTimeout?: string,
  unhealthyConditions?: UnhealthyCondition[],
  status?: MachineHealthCheckStatus,
}

export type Selector = {
  matchLabels?: {[key: string]: string},
  matchExpressions?: {
    key?: string,
    operator?: string,
    values?: string[],
  }[],
}

export type UnhealthyCondition = {
  type?: string,
  status?: string,
  timeout?: string,
}

export type MachineHealthCheckStatus = {
  expectedMachines?: number,
  currentHealthy?: number,
  remediationsAllowed?: number,
  observedGeneration?: number,
  targets?: string[],
  conditions?: Conditions[],
}
