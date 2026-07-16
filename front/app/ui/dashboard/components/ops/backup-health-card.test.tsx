import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {BackupHealthCard} from './backup-health-card';
import {BackupHealth} from '@/app/ui/dashboard/shared/use-day2-ops';

const baseHealth: BackupHealth = {
  available: true,
  storageLocations: [],
  clusterCoverage: [],
  rpoThresholdSeconds: 86400,
  restoresInProgress: 0,
};

describe('BackupHealthCard', () => {
  it('shows a "not available" state when Velero is not installed', () => {
    render(<BackupHealthCard health={{...baseHealth, available: false}}/>);

    expect(screen.getByText(/not available/i)).toBeInTheDocument();
  });

  it('shows storage locations as reachable or unreachable', () => {
    render(<BackupHealthCard health={{
      ...baseHealth,
      storageLocations: [
        {name: 'default', namespace: 'velero', reachable: true, default: true},
        {name: 'secondary', namespace: 'velero', reachable: false, default: false},
      ],
    }}/>);

    expect(screen.getByText('default')).toBeInTheDocument();
    expect(screen.getByText('secondary')).toBeInTheDocument();
    expect(screen.getByRole('img', {name: 'Ready'})).toBeInTheDocument();
    expect(screen.getByRole('img', {name: 'Not ready'})).toBeInTheDocument();
  });

  it('shows a covered, on-time cluster', () => {
    render(<BackupHealthCard health={{
      ...baseHealth,
      clusterCoverage: [{
        clusterRef: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'clusters', namespace: 'default', name: 'capi-workload'},
        covered: true, stale: false, restoreInProgress: false, lastRestoreOutcome: '',
        mostRecentBackupAge: '3h0m0s', mostRecentBackupName: 'nightly-1',
      }],
    }}/>);

    expect(screen.getByText('capi-workload')).toBeInTheDocument();
    expect(screen.getByText(/3h0m0s/)).toBeInTheDocument();
  });

  it('shows a stale cluster distinctly from an on-time one', () => {
    render(<BackupHealthCard health={{
      ...baseHealth,
      clusterCoverage: [{
        clusterRef: {group: '', version: '', resource: '', namespace: 'default', name: 'aging-cluster'},
        covered: true, stale: true, restoreInProgress: false, lastRestoreOutcome: '',
        mostRecentBackupAge: '720h0m0s', mostRecentBackupName: 'old-backup',
      }],
    }}/>);

    expect(screen.getByText(/Stale/)).toBeInTheDocument();
  });

  it('shows a cluster with no backup coverage, not silently omitted', () => {
    render(<BackupHealthCard health={{
      ...baseHealth,
      clusterCoverage: [{
        clusterRef: {group: '', version: '', resource: '', namespace: 'default', name: 'uncovered-cluster'},
        covered: false, stale: false, restoreInProgress: false, lastRestoreOutcome: '',
      }],
    }}/>);

    expect(screen.getByText('uncovered-cluster')).toBeInTheDocument();
    expect(screen.getByText(/no backup/i)).toBeInTheDocument();
  });

  it('shows the restores-in-progress count when greater than zero', () => {
    render(<BackupHealthCard health={{...baseHealth, restoresInProgress: 2}}/>);

    expect(screen.getByText(/2/)).toBeInTheDocument();
    expect(screen.getByText(/restore/i)).toBeInTheDocument();
  });
});
