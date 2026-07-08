import '@testing-library/jest-dom';
import {fireEvent, screen, waitFor} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {DebuggingPath} from './debugging-path';
import {DebugPath} from '@/app/ui/dashboard/shared/use-day2-ops';
import {getDay2OpsDetail, getNodeAccess} from '@/app/lib/data';

jest.mock('@/app/lib/data', () => ({
  getDay2OpsDetail: jest.fn(),
  getNodeAccess: jest.fn(),
}));

const mockedGetDay2OpsDetail = getDay2OpsDetail as jest.MockedFunction<typeof getDay2OpsDetail>;
const mockedGetNodeAccess = getNodeAccess as jest.MockedFunction<typeof getNodeAccess>;

const cappedPath: DebugPath = {
  objectRef: {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'machines', namespace: 'default', name: 'worker-0'},
  layers: [
    {layer: 'conditions', status: 'implicated', evidence: ['Ready=False: WaitingForInfrastructure'], source: 'Machine/worker-0'},
    {layer: 'phase', status: 'implicated', evidence: ['Phase=Provisioning'], source: 'Machine/worker-0'},
    {layer: 'provider_resource', status: 'implicated', evidence: ['VM creation failed'], source: 'DockerMachine/worker-0'},
    {layer: 'controller_activity', status: 'inconclusive', evidence: [], source: ''},
  ],
  summary: 'Waiting on infrastructure provisioning (DockerMachine/worker-0: VM creation failed)',
};

describe('DebuggingPath', () => {
  afterEach(() => jest.resetAllMocks());

  it('renders every layer in order, labeled with its position', () => {
    render(<DebuggingPath path={cappedPath}/>);

    const labels = screen.getAllByText(/^\d$/);
    expect(labels.map((el) => el.textContent)).toEqual(['1', '2', '3', '4']);
    expect(screen.getByText(/Object conditions/)).toBeInTheDocument();
    expect(screen.getByText(/Machine phase/)).toBeInTheDocument();
    expect(screen.getByText(/Provider resource/)).toBeInTheDocument();
    expect(screen.getByText(/Controller activity/)).toBeInTheDocument();
  });

  it('renders the path summary', () => {
    render(<DebuggingPath path={cappedPath}/>);
    expect(screen.getByText(cappedPath.summary)).toBeInTheDocument();
  });

  it('fetches and shows the full evidence list on expand', async () => {
    const fullPath: DebugPath = {
      ...cappedPath,
      layers: cappedPath.layers.map((l) =>
        l.layer === 'conditions' ? {...l, evidence: ['Ready=False: WaitingForInfrastructure', 'InfrastructureReady=False: quota exceeded']} : l,
      ),
    };
    mockedGetDay2OpsDetail.mockResolvedValue({objectRef: cappedPath.objectRef, path: fullPath});

    render(<DebuggingPath path={cappedPath}/>);
    fireEvent.click(screen.getByRole('button', {name: /show full evidence/i}));

    await waitFor(() => {
      expect(screen.getByText(/quota exceeded/)).toBeInTheDocument();
    });
    expect(mockedGetDay2OpsDetail).toHaveBeenCalledWith(
      expect.objectContaining({name: 'worker-0', namespace: 'default'}),
    );
  });

  it('shows an error when the detail fetch fails', async () => {
    mockedGetDay2OpsDetail.mockRejectedValue(new Error('network error'));

    render(<DebuggingPath path={cappedPath}/>);
    fireEvent.click(screen.getByRole('button', {name: /show full evidence/i}));

    await waitFor(() => {
      expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
    });
  });

  it('shows a "View controller logs" deep-dive only when controller_activity is implicated', () => {
    render(<DebuggingPath path={cappedPath}/>);
    expect(screen.queryByRole('link', {name: /view controller logs/i})).not.toBeInTheDocument();

    const implicatedPath: DebugPath = {
      ...cappedPath,
      layers: cappedPath.layers.map((l) =>
        l.layer === 'controller_activity' ? {...l, status: 'implicated', evidence: ['Warning FailedCreate: quota exceeded']} : l,
      ),
    };
    render(<DebuggingPath path={implicatedPath}/>);
    const link = screen.getByRole('link', {name: /view controller logs/i});
    expect(link).toHaveAttribute('href', expect.stringContaining('/dashboard/logs?namespace=capd-system'));
  });

  it('toggles the node-access deep-dive panel', async () => {
    mockedGetNodeAccess.mockResolvedValue({
      objectRef: cappedPath.objectRef,
      command: 'ssh capi@10.0.1.23',
      note: 'Observātiō does not store or manage SSH credentials.',
    });

    render(<DebuggingPath path={cappedPath}/>);
    expect(screen.queryByText(/ssh capi@/)).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', {name: /^node access$/i}));

    await waitFor(() => {
      expect(screen.getByText('ssh capi@10.0.1.23')).toBeInTheDocument();
    });
  });
});
