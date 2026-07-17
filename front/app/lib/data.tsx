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

// ----- AI assistant tool source status (specs/009-mcp-server-client-aggregator) -----

export type MCPSourceKind = 'local' | 'external';
export type MCPHealthState = 'healthy' | 'unhealthy' | 'unknown';

export interface MCPHealthStatus {
  state: MCPHealthState;
  lastChecked?: string;
  lastError?: string;
}

export interface MCPSourceStatus {
  name: string;
  kind: MCPSourceKind;
  health: MCPHealthStatus;
  capabilities: string[];
}

export interface MCPConflict {
  capabilityName: string;
  winningSource: string;
  rejectedSource: string;
}

export interface MCPSourcesResponse {
  sources: MCPSourceStatus[];
  conflicts: MCPConflict[];
}

export const emptyMCPSourcesResponse: MCPSourcesResponse = {sources: [], conflicts: []};

export async function getMCPSources(): Promise<MCPSourcesResponse> {
  return getJSON<MCPSourcesResponse>(`/api/mcp/sources`)
}

// ----- Raw object (YAML tree tab) -----

export interface ResourceGVR {
  group: string;
  version: string;
  resource: string;
}

export async function getRawObject(params: ResourceGVR & { namespace: string; name: string }): Promise<unknown> {
  const query = new URLSearchParams({
    group: params.group,
    version: params.version,
    resource: params.resource,
    namespace: params.namespace,
    name: params.name,
  }).toString();
  return getJSON(`/api/raw?${query}`)
}

// ----- Day-2 Ops debugging-path detail (on-demand drill-in) -----

import type {DebugPath, ObjectRef} from '@/app/ui/dashboard/shared/use-day2-ops';

interface Day2OpsDetailResponse {
  objectRef: ObjectRef;
  path: DebugPath;
}

/** Full, uncapped debugging-path evidence for one object (contracts/day2ops-detail-api.md). */
export async function getDay2OpsDetail(params: ResourceGVR & { namespace: string; name: string }): Promise<Day2OpsDetailResponse> {
  const query = new URLSearchParams({
    group: params.group,
    version: params.version,
    resource: params.resource,
    namespace: params.namespace,
    name: params.name,
  }).toString();
  return getJSON(`/api/day2ops/detail?${query}`)
}

// ----- Day-2 Ops Logs destination (User Story 5) -----

/** Fetches a controller's Pod log output (contracts/logs-api.md). A bounded snapshot, not a live tail. */
export async function getControllerLogs(namespace: string, deployment: string): Promise<string> {
  const query = new URLSearchParams({namespace, deployment}).toString();
  const res = await fetch(`${API_URL}/api/logs/controller?${query}`);
  if (!res.ok) {
    throw new Error(`Request to /api/logs/controller failed: ${res.status} ${res.statusText}`);
  }
  return res.text();
}

export interface NodeAccessInfo {
  objectRef: ObjectRef;
  command: string;
  note: string;
}

/** Static SSH connection instructions for a Machine's node (contracts/logs-api.md). */
export async function getNodeAccess(params: ResourceGVR & { namespace: string; name: string }): Promise<NodeAccessInfo> {
  const query = new URLSearchParams({
    group: params.group,
    version: params.version,
    resource: params.resource,
    namespace: params.namespace,
    name: params.name,
  }).toString();
  return getJSON(`/api/logs/node-access?${query}`)
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
