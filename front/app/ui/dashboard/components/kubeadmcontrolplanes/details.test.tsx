import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {Grid} from '@mantine/core';
import {render as baseRender} from '@/app/ui/dashboard/utils/test-render';
import {AIPanelProvider} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

import KubeadmControlPlaneDetails from './details';
import {KubeadmControlPlaneType} from './types';

const render = (ui: React.ReactNode) => baseRender(<AIPanelProvider><Grid>{ui}</Grid></AIPanelProvider>);

const kcp: KubeadmControlPlaneType = {
  metadata: {name: 'mgmt-kcp', namespace: 'default'},
  age: '10d',
  cluster: 'capi-mgmt',
  version: 'v1.31.0',
  replicas: 3,
  status: {
    readyReplicas: 3, updatedReplicas: 3, unavailableReplicas: 0, initialized: true, ready: true,
    conditions: [{type: 'EtcdClusterHealthy', status: 'True'}],
  },
};

describe('KubeadmControlPlaneDetails', () => {
  it('shows Specification and YAML tabs, no embedded AI Troubleshooting tab', () => {
    render(<KubeadmControlPlaneDetails kcp={kcp}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
    expect(screen.getByRole('tab', {name: 'YAML'})).toBeInTheDocument();
    expect(screen.queryByRole('tab', {name: 'AI Troubleshooting'})).not.toBeInTheDocument();
  });

  it('renders desired/ready replicas, version, and etcd conditions', () => {
    render(<KubeadmControlPlaneDetails kcp={kcp}/>);
    expect(screen.getByText('v1.31.0')).toBeInTheDocument();
    expect(screen.getByText('capi-mgmt')).toBeInTheDocument();
    expect(screen.getByText('EtcdClusterHealthy')).toBeInTheDocument();
  });

  it('renders a partial KubeadmControlPlane (missing status) without throwing', () => {
    const partial = {metadata: {name: 'kcp1'}} as KubeadmControlPlaneType;
    render(<KubeadmControlPlaneDetails kcp={partial}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
  });
});
