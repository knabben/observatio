'use client';

import {useCallback, useEffect, useState} from 'react';
import {sendInitialRequest, WebSocket} from '@/app/lib/websocket';

export interface ObjectRef {
  group: string;
  version: string;
  resource: string;
  namespace: string;
  name: string;
}

export type Category = 'cluster' | 'machine_deployment' | 'machine';

export interface HealthRollup {
  category: Category;
  healthy: number;
  degraded: number;
  failed: number;
  unavailable: boolean;
}

export type DebugLayerName = 'conditions' | 'phase' | 'provider_resource' | 'controller_activity';
export type DebugLayerStatus = 'ok' | 'implicated' | 'inconclusive';

export interface DebugLayer {
  layer: DebugLayerName;
  status: DebugLayerStatus;
  evidence: string[];
  source: string;
}

export interface DebugPath {
  objectRef: ObjectRef;
  layers: DebugLayer[];
  summary: string;
}

export type RiskKind = 'cert_expiry' | 'stalled_rollout' | 'version_skew' | 'drift';
export type RiskCheckStatus = 'evaluated' | 'not_evaluable';

export interface RiskWarning {
  objectRef: ObjectRef;
  kind: RiskKind;
  detail: string;
  likelyCause: string;
  checkStatus: RiskCheckStatus;
}

export type SeverityLevel =
  | 'self_healing'
  | 'needs_investigation'
  | 'provider_degraded'
  | 'management_critical';

export interface RecoveryInfo {
  recoverable: boolean;
  coveringBackupAge?: string;
}

export interface FailureSeverity {
  objectRef: ObjectRef | null;
  level: SeverityLevel;
  reason: string;
  recoveryInfo?: RecoveryInfo;
}

export interface BackupStorageLocationStatus {
  name: string;
  namespace: string;
  reachable: boolean;
  default: boolean;
}

export type LastRestoreOutcome = '' | 'succeeded' | 'failed';

export interface ClusterBackupCoverage {
  clusterRef: ObjectRef;
  covered: boolean;
  mostRecentBackupAge?: string;
  mostRecentBackupName?: string;
  stale: boolean;
  restoreInProgress: boolean;
  lastRestoreOutcome: LastRestoreOutcome;
}

export interface BackupHealth {
  available: boolean;
  storageLocations: BackupStorageLocationStatus[];
  clusterCoverage: ClusterBackupCoverage[];
  rpoThresholdSeconds: number;
  restoresInProgress: number;
}

export interface Day2OpsData {
  rollups: HealthRollup[];
  debugPaths: DebugPath[];
  risks: RiskWarning[];
  severities: FailureSeverity[];
  sourceUnavailable: boolean;
  backupHealth: BackupHealth;
}

const EMPTY_BACKUP_HEALTH: BackupHealth = {
  available: false,
  storageLocations: [],
  clusterCoverage: [],
  rpoThresholdSeconds: 0,
  restoresInProgress: 0,
};

const EMPTY_DATA: Day2OpsData = {
  rollups: [],
  debugPaths: [],
  risks: [],
  severities: [],
  sourceUnavailable: false,
  backupHealth: EMPTY_BACKUP_HEALTH,
};

interface Day2OpsState {
  data: Day2OpsData;
  /** True once the first frame has arrived; distinguishes "still connecting" from "all clear". */
  loaded: boolean;
}

/**
 * Subscribes to the `day2ops` WebSocket event (contracts/day2ops-ws-event.md). Every frame is a
 * full-state replace, not a patch — unlike `useResourceStream`, there is no per-item merge here.
 */
export function useDay2Ops(): Day2OpsState {
  const [data, setData] = useState<Day2OpsData>(EMPTY_DATA);
  const [loaded, setLoaded] = useState(false);

  const onReconnectStop = useCallback(() => {
    setData((prev) => ({...prev, sourceUnavailable: true}));
  }, []);

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket(undefined, {onReconnectStop});

  useEffect(() => {
    sendInitialRequest(readyState, 'day2ops', sendJsonMessage);
  }, [readyState, sendJsonMessage]);

  useEffect(() => {
    const message = lastJsonMessage as {event?: string; data?: Day2OpsData} | null;
    if (message?.event !== 'day2ops' || message.data == null) return;
    setData(message.data);
    setLoaded(true);
  }, [lastJsonMessage]);

  return {data, loaded};
}
