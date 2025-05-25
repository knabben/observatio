import {Conditions, Meta} from "@/app/ui/dashboard/base/types";

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
  metadata: Meta,
  replicas: number,
  cluster: string,
  age: string,
  templateversion: string,
  templateBootstrap: Bootstrap,
  templateInfrastructureRef: InfrastructureRef,
  status: MachineDeploymentStatus,
}

/**
 * Represents a Bootstrap type definition that contains a reference
 * configuration for a resource.
 */
export type Bootstrap = {
  configRef: {
    name: string,
    namespace: string,
    kind: string,
    apiVersion: string,
  }
}

/**
 * Represents a reference to an infrastructure resource within a Kubernetes-like environment.
 *
 * This type is used to uniquely identify a resource by specifying its name, namespace, kind,
 * and API version. The `InfrastructureRef` structure allows for proper cross-referencing
 * of resources in a cluster or a system that adheres to Kubernetes conventions.
 *
 * Properties:
 * - `name`: The name of the resource being referenced.
 * - `namespace`: The namespace in which the resource resides (if applicable).
 * - `kind`: The specific resource type or kind (e.g., Deployment, Service, etc.).
 * - `apiVersion`: The API version of the resource being referenced.
 */
export type InfrastructureRef = {
  name: string,
  namespace: string,
  kind: string,
  apiVersion: string,
}

/**
 * Represents the status of a machine deployment within a cluster.
 *
 * This type defines the current state of the deployment, including information
 * about its replicas, conditions, and deployment phase.
 *
 * Properties:
 * - `phase`: Describes the phase of the deployment (e.g., "Pending", "Running", "Failed").
 * - `conditions`: An array of conditions that provide additional details about the deployment's current state.
 * - `readyReplicas`: The number of replicas that are currently ready and available.
 * - `updatedReplicas`: The number of replicas that have been updated to the desired specification.
 * - `unavailableReplicas`: The number of replicas that are unavailable and not ready.
 * - `v1beta2`: Contains additional deployment status details, including:
 *   - `availableReplicas`: Number of replicas that are actively available.
 *   - `readyReplicas`: Number of replicas that are fully configured and ready.
 *   - `upToDateReplicas`: Number of replicas that are up-to-date with the desired configuration.
 *   - `conditions`: An array of status conditions specific to v1beta2 upgrades.
 */
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