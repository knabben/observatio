import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {Grid} from '@mantine/core';
import {render as baseRender} from '@/app/ui/dashboard/utils/test-render';
import {AIPanelProvider} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

import ClusterClassDetails from './details';
import {ClusterClassType} from './types';

const render = (ui: React.ReactNode) => baseRender(<AIPanelProvider><Grid>{ui}</Grid></AIPanelProvider>);

const cc: ClusterClassType = {
  metadata: {name: 'quick-start', namespace: 'default'},
  name: 'quick-start',
  namespace: 'default',
  generation: 2,
  conditions: [{type: 'RefVersionsUpToDate', status: 'True'}],
};

describe('ClusterClassDetails', () => {
  it('shows Specification and YAML tabs, no embedded AI Troubleshooting tab', () => {
    render(<ClusterClassDetails cc={cc}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
    expect(screen.getByRole('tab', {name: 'YAML'})).toBeInTheDocument();
    expect(screen.queryByRole('tab', {name: 'AI Troubleshooting'})).not.toBeInTheDocument();
  });

  it('renders status/reference fields', () => {
    render(<ClusterClassDetails cc={cc}/>);
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('RefVersionsUpToDate')).toBeInTheDocument();
  });

  it('renders a partial ClusterClass (missing conditions) without throwing', () => {
    const partial = {metadata: {name: 'cc1'}, name: 'cc1'} as ClusterClassType;
    render(<ClusterClassDetails cc={partial}/>);
    expect(screen.getByRole('tab', {name: 'Specification'})).toBeInTheDocument();
  });
});
