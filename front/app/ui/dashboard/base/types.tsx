
/**
 * Represents metadata information about a resource, often used in Kubernetes and similar systems.
 * Includes details such as name, namespace, and timestamps, along with additional labels, annotations,
 * and ownership references.
 *
 * Properties:
 * - `name`: The name of the resource.
 * - `namespace`: The namespace the resource belongs to.
 * - `resourceVersion`: A unique identifier representing the version of the resource for consistency tracking.
 * - `creationTimestamp`: The timestamp of when the resource was created.
 * - `labels`: A map of key-value pairs used to organize and categorize resources.
 * - `annotations`: A map of key-value pairs for storing arbitrary metadata related to the resource.
 * - `ownerReferences`: An array of references to owner resources, describing the relationship between resources.
 */
export type Meta = {
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

/**
 * A type representing the state or conditions associated with an object.
 */
export type Conditions = {
  type: string,
  status: string,
  severity: string,
  lastTransitionTime: string,
  reason: string,
  message: string,
}

