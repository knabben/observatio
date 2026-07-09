import {Conditions, Meta} from "@/app/ui/dashboard/base/types";

/**
 * Represents the type definition for a KubeadmControlPlane — control-plane replica health and
 * status conditions (including etcd-related conditions when present).
 */
export type KubeadmControlPlaneType = {
  metadata?: Meta,
  age?: string,
  cluster?: string,
  version?: string,
  replicas?: number,
  status?: KubeadmControlPlaneStatus,
}

export type KubeadmControlPlaneStatus = {
  selector?: string,
  replicas?: number,
  version?: string,
  updatedReplicas?: number,
  readyReplicas?: number,
  unavailableReplicas?: number,
  initialized?: boolean,
  ready?: boolean,
  observedGeneration?: number,
  conditions?: Conditions[],
}
