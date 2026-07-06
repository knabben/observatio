/* eslint-disable @typescript-eslint/no-explicit-any */
import {API_URL} from "@/app/lib/config";

/**
 * Fetches JSON from the backend, checking the HTTP status first so a 4xx/5xx body is
 * never silently treated as valid data. Network and parse failures propagate to the
 * caller, which surfaces them as an error state.
 */
async function getJSON<T = any>(path: string): Promise<T> {
  const res = await fetch(`${API_URL}${path}`);
  if (!res.ok) {
    throw new Error(`Request to ${path} failed: ${res.status} ${res.statusText}`);
  }
  return res.json() as Promise<T>;
}

// ----- Infrastructure provider detection -----

export interface ProviderStatus {
  installed: boolean;
  version: string;
}

export interface InfrastructureCapability {
  docker: ProviderStatus;
  vsphere: ProviderStatus;
}

export const emptyInfrastructureCapability: InfrastructureCapability = {
  docker: {installed: false, version: ''},
  vsphere: {installed: false, version: ''},
};

export async function getInfraCapabilities(): Promise<InfrastructureCapability> {
  return getJSON<InfrastructureCapability>(`/api/infra/capabilities`)
}

// ----- Dashboard -----

export async function getComponentsVersion() {
  return getJSON(`/api/clusters/components`)
}

export async function getClusterInformation() {
  return getJSON(`/api/clusters/info`)
}

export async function getClusterSummary() {
  return getJSON(`/api/clusters/summary`)
}

export async function getClusterClasses() {
  return getJSON(`/api/clusters/classes`)
}

export async function getClusterHierarchy() {
  return getJSON(`/api/clusters/topology`)
}

// ----- Clusters -----

export async function getClusterList() {
  return getJSON(`/api/clusters/list`)
}

export async function getClusterInfraList() {
  return getJSON(`/api/clusters/infra/list`)
}

// ----- MachinesDeployment -----

export async function getMachinesDeployments() {
  return getJSON(`/api/machinesdeployment/list`)
}

// ----- Machines -----

export async function getMachines() {
  return getJSON(`/api/machines/list`)
}

// ----- AI Interactions -----
const defaultRequestConfig = {
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
};

export async function postAIAnalysis(request: string) {
  const res = await fetch(`${API_URL}/api/analysis`, {
    ...defaultRequestConfig,
    method: 'POST',
    body: JSON.stringify({request}),
  });
  if (!res.ok) {
    throw new Error(`AI analysis request failed: ${res.status} ${res.statusText}`);
  }
  return res.json();
}
