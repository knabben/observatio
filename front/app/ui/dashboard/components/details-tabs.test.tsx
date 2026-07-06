import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {Grid} from '@mantine/core';
import {render as baseRender} from '@/app/ui/dashboard/utils/test-render';
import {AIPanelProvider} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

const render = (ui: React.ReactNode) => baseRender(<AIPanelProvider><Grid>{ui}</Grid></AIPanelProvider>);

import ClusterDetails from './clusters/details';
import ClusterInfraDetails from './clusters/infra/infra-details';
import ClusterInfraDockerDetails from './clusters/infra/docker-details';
import MachineDetails from './machines/details';
import MachineInfraDetails from './machines/infra/infra-details';
import MachineInfraDockerDetails from './machines/infra/docker-details';
import MachineDeploymentDetails from './mds/details';

const meta = {name: 'r1', namespace: 'default'};

function expectSpecificationTabOnlyNoAI() {
  expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
  expect(screen.queryByRole('tab', {name: 'AI Troubleshooting'})).not.toBeInTheDocument();
}

describe('Object detail screens: no embedded AI Troubleshooting tab, Specification + conditions remain', () => {
  it('ClusterDetails', () => {
    render(<ClusterDetails cluster={{
      metadata: meta,
      status: {conditions: [{type: 'Ready', status: 'True'}]},
    }}/>);
    expectSpecificationTabOnlyNoAI();
    expect(screen.getByText('Ready')).toBeInTheDocument();
  });

  it('ClusterInfraDetails (vSphere)', () => {
    render(<ClusterInfraDetails cluster={{
      metadata: meta,
      status: {conditions: [{type: 'Ready', status: 'True'}]},
    }}/>);
    expectSpecificationTabOnlyNoAI();
    expect(screen.getByText('Ready')).toBeInTheDocument();
  });

  it('ClusterInfraDockerDetails', () => {
    render(<ClusterInfraDockerDetails cluster={{metadata: meta, ready: true}}/>);
    expectSpecificationTabOnlyNoAI();
  });

  it('MachineDetails', () => {
    render(<MachineDetails machine={{
      metadata: meta,
      status: {conditions: [{type: 'Ready', status: 'True'}]},
    }}/>);
    expectSpecificationTabOnlyNoAI();
    expect(screen.getByText('Ready')).toBeInTheDocument();
  });

  it('MachineInfraDetails (vSphere)', () => {
    render(<MachineInfraDetails machine={{
      metadata: meta,
      status: {conditions: [{type: 'Ready', status: 'True'}]},
    }}/>);
    expectSpecificationTabOnlyNoAI();
    expect(screen.getByText('Ready')).toBeInTheDocument();
  });

  it('MachineInfraDockerDetails', () => {
    render(<MachineInfraDockerDetails machine={{metadata: meta, ready: true}}/>);
    expectSpecificationTabOnlyNoAI();
  });

  it('MachineDeploymentDetails', () => {
    render(<MachineDeploymentDetails md={{
      metadata: meta,
      status: {conditions: [{type: 'Ready', status: 'True'}]},
    }}/>);
    expectSpecificationTabOnlyNoAI();
    expect(screen.getByText('Ready')).toBeInTheDocument();
  });
});
