// MachineDeployment models and details

type MachineDeploymentType = {
  name: string,
  replicas: number,
  namespace: string,
  cluster: string,
  readyReplicas: number,
  updatedReplicas: number,
  unavailableReplicas: number,
  created: string
  phase: string,
}

