import '@testing-library/jest-dom';
import {screen, waitFor} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {NodeAccessPanel} from './node-access-panel';
import {getNodeAccess} from '@/app/lib/data';

jest.mock('@/app/lib/data', () => ({
  getNodeAccess: jest.fn(),
}));

const mockedGetNodeAccess = getNodeAccess as jest.MockedFunction<typeof getNodeAccess>;

const objectRef = {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'machines', namespace: 'default', name: 'worker-0'};

describe('NodeAccessPanel', () => {
  afterEach(() => jest.resetAllMocks());

  it('renders only the SSH command, address, and disclaimer — no credential input', async () => {
    mockedGetNodeAccess.mockResolvedValue({
      objectRef,
      command: 'ssh capi@10.0.1.23',
      note: 'Observātiō does not store or manage SSH credentials. Run this command from your own machine.',
    });

    const {container} = render(<NodeAccessPanel objectRef={objectRef}/>);

    await waitFor(() => {
      expect(screen.getByText('ssh capi@10.0.1.23')).toBeInTheDocument();
    });
    expect(screen.getByText(/does not store or manage SSH credentials/i)).toBeInTheDocument();
    expect(container.querySelector('input[type="password"]')).toBeNull();
    expect(container.querySelector('input')).toBeNull();
  });

  it('shows an error state when the fetch fails', async () => {
    mockedGetNodeAccess.mockRejectedValue(new Error('404'));

    render(<NodeAccessPanel objectRef={objectRef}/>);

    await waitFor(() => {
      expect(screen.getByText(/could not be determined/i)).toBeInTheDocument();
    });
  });
});
