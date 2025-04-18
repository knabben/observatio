const URL = "http://localhost:8080"

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

// ----- Clusters -----

export async function getClusterList() {
  const res = await fetch(`${URL}/api/clusters/list`)
  return res.json()
}