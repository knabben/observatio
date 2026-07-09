import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {Grid} from '@mantine/core';
import {render as baseRender} from '@/app/ui/dashboard/utils/test-render';
import {AIPanelProvider} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

import MachineSetDetails from './details';
import {MachineSetType} from './types';

const render = (ui: React.ReactNode) => baseRender(<AIPanelProvider><Grid>{ui}</Grid></AIPanelProvider>);

const ms: MachineSetType = {
  metadata: {name: 'worker-ms', namespace: 'default'},
  age: '5d',
  cluster: 'capi-workload',
  machineDeployment: 'worker-md',
  replicas: 3,
  status: {readyReplicas: 3, availableReplicas: 3, fullyLabeledReplicas: 3},
};

describe('MachineSetDetails', () => {
  it('shows Specification and YAML tabs, no embedded AI Troubleshooting tab', () => {
    render(<MachineSetDetails ms={ms}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
    expect(screen.getByRole('tab', {name: 'YAML'})).toBeInTheDocument();
    expect(screen.queryByRole('tab', {name: 'AI Troubleshooting'})).not.toBeInTheDocument();
  });

  it('renders replica counts and the owning MachineDeployment', () => {
    render(<MachineSetDetails ms={ms}/>);
    expect(screen.getByText('worker-md')).toBeInTheDocument();
    expect(screen.getByText('capi-workload')).toBeInTheDocument();
  });

  it('renders a partial MachineSet (missing status) without throwing', () => {
    const partial = {metadata: {name: 'ms1'}} as MachineSetType;
    render(<MachineSetDetails ms={partial}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
  });
});
