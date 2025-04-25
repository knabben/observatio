// Cluster types and models.

type Conditions = {
  type: string,
  status: boolean,
  lastTransitionTime: string,
}

export type ClusterType = {
  name: string,
  hasTopology: boolean,
  conditions: Conditions[]
}

