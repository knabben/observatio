import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {Grid} from '@mantine/core';
import {render as baseRender} from '@/app/ui/dashboard/utils/test-render';
import {AIPanelProvider} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

import MachineHealthCheckDetails from './details';
import {MachineHealthCheckType} from './types';

const render = (ui: React.ReactNode) => baseRender(<AIPanelProvider><Grid>{ui}</Grid></AIPanelProvider>);

const mhc: MachineHealthCheckType = {
  metadata: {name: 'worker-mhc', namespace: 'default'},
  age: '3d',
  cluster: 'capi-workload',
  selector: {matchLabels: {role: 'worker'}},
  maxUnhealthy: '40%',
  nodeStartupTimeout: '10m0s',
  unhealthyConditions: [{type: 'Ready', status: 'False', timeout: '5m0s'}],
  status: {expectedMachines: 3, currentHealthy: 2, remediationsAllowed: 1, conditions: [{type: 'RemediationAllowed', status: 'True'}]},
};

describe('MachineHealthCheckDetails', () => {
  it('shows Specification and YAML tabs, no embedded AI Troubleshooting tab', () => {
    render(<MachineHealthCheckDetails mhc={mhc}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
    expect(screen.getByRole('tab', {name: 'YAML'})).toBeInTheDocument();
    expect(screen.queryByRole('tab', {name: 'AI Troubleshooting'})).not.toBeInTheDocument();
  });

  it('renders the target selector, maxUnhealthy, timeouts, and remediation status', () => {
    render(<MachineHealthCheckDetails mhc={mhc}/>);
    expect(screen.getByText('role=worker')).toBeInTheDocument();
    expect(screen.getByText('40%')).toBeInTheDocument();
    expect(screen.getByText('10m0s')).toBeInTheDocument();
    expect(screen.getByText('Ready=False (5m0s)')).toBeInTheDocument();
    expect(screen.getByText('capi-workload')).toBeInTheDocument();
  });

  it('renders a partial MachineHealthCheck (missing status) without throwing', () => {
    const partial = {metadata: {name: 'mhc1'}} as MachineHealthCheckType;
    render(<MachineHealthCheckDetails mhc={partial}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
  });
});
