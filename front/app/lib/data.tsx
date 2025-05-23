export const URL = 'http://localhost:8080';

// ----- Dashboard -----

export async function getComponentsVersion() {
  const res = await fetch(`${URL}/api/clusters/components`)
  return res.json()
}

export async function getClusterInformation() {
  const res = await fetch(`${URL}/api/clusters/info`)
  return res.json()
}

export async function getClusterSummary() {
  const res = await fetch(`${URL}/api/clusters/summary`)
  return res.json()
}

export async function getClusterClasses() {
  const res = await fetch(`${URL}/api/clusters/classes`)
  return res.json()
}

export async function getClusterHierarchy() {
  const res = await fetch(`${URL}/api/clusters/topology`)
  return res.json()
}

// ----- Clusters -----

export async function getClusterList() {
  const res = await fetch(`${URL}/api/clusters/list`)
  return res.json()
}

export async function getClusterInfraList() {
  const res = await fetch(`${URL}/api/clusters/infra/list`)
  return res.json()
}

// ----- MachinesDeployment -----

export async function getMachinesDeployments() {
  const res = await fetch(`${URL}/api/machinesdeployment/list`)
  return res.json()
}

// ----- Machines -----

export async function getMachines() {
  const res = await fetch(`${URL}/api/machines/list`)
  return res.json()
}

// ----- AI Interactions -----
const defaultRequestConfig = {
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
};

export async function postAIAnalysis(request: string) {
  const res = await fetch(`${URL}/api/analysis`, {
    ...defaultRequestConfig,
    method: 'POST',
    body: JSON.stringify({request}),
  });
  return res.json();
}
