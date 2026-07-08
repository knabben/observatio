import '@testing-library/jest-dom';
import {fireEvent, screen, waitFor} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {LogsView} from './logs-view';
import {getControllerLogs} from '@/app/lib/data';

jest.mock('@/app/lib/data', () => ({
  getControllerLogs: jest.fn(),
}));

let mockSearchParams = new URLSearchParams();
jest.mock('next/navigation', () => ({
  useSearchParams: () => mockSearchParams,
}));

const mockedGetControllerLogs = getControllerLogs as jest.MockedFunction<typeof getControllerLogs>;

beforeAll(() => {
  (global as unknown as {ResizeObserver: unknown}).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
});

describe('LogsView', () => {
  afterEach(() => {
    jest.resetAllMocks();
    mockSearchParams = new URLSearchParams();
  });

  it('preselects a controller from the URL query params (deep-dive link)', async () => {
    mockedGetControllerLogs.mockResolvedValue('capd log from query\n');
    mockSearchParams = new URLSearchParams({namespace: 'capd-system', deployment: 'capd-controller-manager'});

    render(<LogsView/>);

    await waitFor(() => {
      expect(screen.getByText(/capd log from query/)).toBeInTheDocument();
    });
  });

  it('lets the operator choose a controller and shows its log output', async () => {
    mockedGetControllerLogs.mockResolvedValue('line 1\nline 2\n');

    render(<LogsView/>);
    fireEvent.click(screen.getByRole('button', {name: /capi-controller-manager/i}));

    await waitFor(() => {
      expect(screen.getByText(/line 1/)).toBeInTheDocument();
    });
    expect(mockedGetControllerLogs).toHaveBeenCalledWith('capi-system', 'capi-controller-manager');
  });

  it('preselects a controller when initialController is provided', async () => {
    mockedGetControllerLogs.mockResolvedValue('capd log line\n');

    render(<LogsView initialController={{namespace: 'capd-system', deployment: 'capd-controller-manager'}}/>);

    await waitFor(() => {
      expect(screen.getByText(/capd log line/)).toBeInTheDocument();
    });
  });

  it('shows a "logs unavailable" state when the fetch fails', async () => {
    mockedGetControllerLogs.mockRejectedValue(new Error('503'));

    render(<LogsView initialController={{namespace: 'capi-system', deployment: 'capi-controller-manager'}}/>);

    await waitFor(() => {
      expect(screen.getByText(/logs unavailable/i)).toBeInTheDocument();
    });
  });
});
