import '@testing-library/jest-dom';
import {screen, waitFor} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';

jest.mock('@/app/lib/data');

import {getRawObject} from '@/app/lib/data';
import {ObjectTree} from './object-tree';

const mockedGetRawObject = getRawObject as jest.MockedFunction<typeof getRawObject>;

beforeAll(() => {
  (global as unknown as {ResizeObserver: unknown}).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
});

afterEach(() => jest.resetAllMocks());

const gvr = {group: 'cluster.x-k8s.io', version: 'v1beta1', resource: 'clusters'};

describe('ObjectTree', () => {
  it('renders the fetched object as a tree', async () => {
    mockedGetRawObject.mockResolvedValue({metadata: {name: 'c1'}, spec: {paused: false}});
    render(<ObjectTree gvr={gvr} namespace="default" name="c1"/>);
    expect(await screen.findByText('metadata')).toBeInTheDocument();
    expect(screen.getByText('spec')).toBeInTheDocument();
  });

  it('shows an error state when the fetch fails, not a crash', async () => {
    mockedGetRawObject.mockRejectedValue(new Error('network down'));
    render(<ObjectTree gvr={gvr} namespace="default" name="c1"/>);
    expect(await screen.findByText(/failed to load the complete object/i)).toBeInTheDocument();
  });

  it('re-fetches when resourceVersion changes', async () => {
    mockedGetRawObject.mockResolvedValue({metadata: {name: 'c1'}});
    const {rerender} = render(<ObjectTree gvr={gvr} namespace="default" name="c1" resourceVersion="1"/>);
    await waitFor(() => expect(mockedGetRawObject).toHaveBeenCalledTimes(1));

    rerender(<ObjectTree gvr={gvr} namespace="default" name="c1" resourceVersion="2"/>);
    await waitFor(() => expect(mockedGetRawObject).toHaveBeenCalledTimes(2));
  });
});
