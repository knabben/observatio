import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {RiskWarnings} from './risk-warnings';
import {RiskWarning} from '@/app/ui/dashboard/shared/use-day2-ops';

const certExpiry: RiskWarning = {
  objectRef: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'clusters', namespace: 'default', name: 'prod-1'},
  kind: 'cert_expiry',
  detail: 'prod-1-ca expires 2026-08-01',
  likelyCause: '',
  checkStatus: 'evaluated',
};

const stalledRollout: RiskWarning = {
  objectRef: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'machinedeployments', namespace: 'default', name: 'workers'},
  kind: 'stalled_rollout',
  detail: 'MachineSet workers-old has not scaled down after 45m',
  likelyCause: 'blocked by finalizer(s)',
  checkStatus: 'evaluated',
};

const notEvaluable: RiskWarning = {
  objectRef: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'clusters', namespace: 'default', name: 'prod-2'},
  kind: 'cert_expiry',
  detail: '',
  likelyCause: '',
  checkStatus: 'not_evaluable',
};

describe('RiskWarnings', () => {
  it('renders warnings grouped by kind', () => {
    render(<RiskWarnings risks={[certExpiry, stalledRollout]}/>);

    expect(screen.getByText(/certificate expiry/i)).toBeInTheDocument();
    expect(screen.getByText(/stalled rollout/i)).toBeInTheDocument();
    expect(screen.getByText('prod-1-ca expires 2026-08-01')).toBeInTheDocument();
    expect(screen.getByText(/MachineSet workers-old/)).toBeInTheDocument();
  });

  it('shows the likely cause when determinable', () => {
    render(<RiskWarnings risks={[stalledRollout]}/>);
    expect(screen.getByText(/blocked by finalizer/i)).toBeInTheDocument();
  });

  it('shows a "check could not be performed" state instead of omitting a not-evaluable risk', () => {
    render(<RiskWarnings risks={[notEvaluable]}/>);
    expect(screen.getByText(/check could not be performed/i)).toBeInTheDocument();
  });

  it('renders nothing when there are no risks', () => {
    render(<RiskWarnings risks={[]}/>);
    expect(screen.queryByText(/proactive risk warnings/i)).not.toBeInTheDocument();
  });
});
