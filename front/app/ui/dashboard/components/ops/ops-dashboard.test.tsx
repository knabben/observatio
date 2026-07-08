import '@testing-library/jest-dom';
import {fireEvent, screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import useWebSocket, {ReadyState} from 'react-use-websocket';
import {OpsDashboard} from './ops-dashboard';
import {Day2OpsData} from '@/app/ui/dashboard/shared/use-day2-ops';

jest.mock('react-use-websocket');

const mockedUseWebSocket = useWebSocket as jest.MockedFunction<typeof useWebSocket>;

const allHealthyData: Day2OpsData = {
  rollups: [
    {category: 'cluster', healthy: 2, degraded: 0, failed: 0, unavailable: false},
    {category: 'machine_deployment', healthy: 1, degraded: 0, failed: 0, unavailable: false},
    {category: 'machine', healthy: 5, degraded: 0, failed: 0, unavailable: false},
  ],
  debugPaths: [],
  risks: [],
  severities: [],
  sourceUnavailable: false,
};

function mockSocket(lastJsonMessage: unknown, readyState: ReadyState = ReadyState.OPEN) {
  mockedUseWebSocket.mockReturnValue({
    sendJsonMessage: jest.fn(),
    lastJsonMessage,
    readyState,
  } as unknown as ReturnType<typeof useWebSocket>);
}

describe('OpsDashboard', () => {
  afterEach(() => jest.resetAllMocks());

  it('renders every category rollup card', () => {
    mockSocket({type: 'MODIFIED', event: 'day2ops', data: allHealthyData});
    render(<OpsDashboard/>);

    expect(screen.getByRole('heading', {name: 'Clusters'})).toBeInTheDocument();
    expect(screen.getByRole('heading', {name: 'Machine Deployments'})).toBeInTheDocument();
    expect(screen.getByRole('heading', {name: 'Machines'})).toBeInTheDocument();
  });

  it('narrows to a single category in place without navigating away', () => {
    mockSocket({type: 'MODIFIED', event: 'day2ops', data: allHealthyData});
    render(<OpsDashboard/>);

    fireEvent.click(screen.getByRole('button', {name: 'Machines'}));

    expect(screen.getByRole('heading', {name: 'Machines'})).toBeInTheDocument();
    expect(screen.queryByRole('heading', {name: 'Clusters'})).not.toBeInTheDocument();
    expect(screen.queryByRole('heading', {name: 'Machine Deployments'})).not.toBeInTheDocument();
  });

  it('shows a data-unavailable banner instead of a false all-clear state', () => {
    mockSocket({
      type: 'MODIFIED',
      event: 'day2ops',
      data: {...allHealthyData, sourceUnavailable: true},
    });
    render(<OpsDashboard/>);

    expect(screen.getByText(/data unavailable/i)).toBeInTheDocument();
  });

  it('renders the debugging path inline for each unhealthy machine (FR-004)', () => {
    mockSocket({
      type: 'MODIFIED',
      event: 'day2ops',
      data: {
        ...allHealthyData,
        debugPaths: [{
          objectRef: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'machines', namespace: 'default', name: 'worker-0'},
          layers: [
            {layer: 'conditions', status: 'implicated', evidence: ['Ready=False: WaitingForInfrastructure'], source: 'Machine/worker-0'},
            {layer: 'phase', status: 'implicated', evidence: ['Phase=Provisioning'], source: 'Machine/worker-0'},
            {layer: 'provider_resource', status: 'inconclusive', evidence: [], source: ''},
            {layer: 'controller_activity', status: 'inconclusive', evidence: [], source: ''},
          ],
          summary: 'Waiting on object conditions (Machine/worker-0: Ready=False: WaitingForInfrastructure)',
        }],
      },
    });
    render(<OpsDashboard/>);

    expect(screen.getByText(/Waiting on object conditions/)).toBeInTheDocument();
  });
});
