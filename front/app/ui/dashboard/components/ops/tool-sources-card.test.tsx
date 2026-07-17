import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';

jest.mock('@/app/lib/data');

import {getMCPSources, MCPSourcesResponse} from '@/app/lib/data';
import {ToolSourcesCard} from './tool-sources-card';

const mockedGetMCPSources = getMCPSources as jest.MockedFunction<typeof getMCPSources>;

afterEach(() => {
  jest.resetAllMocks();
});

const baseResponse: MCPSourcesResponse = {
  sources: [],
  conflicts: [],
};

describe('ToolSourcesCard', () => {
  it('shows a healthy source with its capabilities', async () => {
    mockedGetMCPSources.mockResolvedValue({
      ...baseResponse,
      sources: [{
        name: 'kubectl', kind: 'local', capabilities: ['kubectl'],
        health: {state: 'healthy'},
      }],
    });

    render(<ToolSourcesCard/>);

    // "kubectl" appears twice — once as the source name, once as its one capability's name.
    expect(await screen.findAllByText('kubectl')).toHaveLength(2);
    expect(screen.getByText('local')).toBeInTheDocument();
    expect(screen.getByRole('img', {name: 'Ready'})).toBeInTheDocument();
  });

  it('shows an unhealthy source distinctly, still listed', async () => {
    mockedGetMCPSources.mockResolvedValue({
      ...baseResponse,
      sources: [{
        name: 'velero-mcp', kind: 'external', capabilities: ['list_backups'],
        health: {state: 'unhealthy', lastError: 'connection refused'},
      }],
    });

    render(<ToolSourcesCard/>);

    expect(await screen.findByText('velero-mcp')).toBeInTheDocument();
    expect(screen.getByRole('img', {name: 'Not ready'})).toBeInTheDocument();
  });

  it('shows an unknown-health source (never yet probed) distinctly from healthy/unhealthy', async () => {
    mockedGetMCPSources.mockResolvedValue({
      ...baseResponse,
      sources: [{
        name: 'brand-new', kind: 'external', capabilities: [],
        health: {state: 'unknown'},
      }],
    });

    render(<ToolSourcesCard/>);

    expect(await screen.findByText('brand-new')).toBeInTheDocument();
    expect(screen.getByRole('img', {name: 'Unknown'})).toBeInTheDocument();
  });

  it('flags a capability a source lost to a naming conflict', async () => {
    mockedGetMCPSources.mockResolvedValue({
      sources: [{
        name: 'velero-mcp-mirror', kind: 'external', capabilities: [],
        health: {state: 'healthy'},
      }],
      conflicts: [{capabilityName: 'list_backups', winningSource: 'velero-mcp', rejectedSource: 'velero-mcp-mirror'}],
    });

    render(<ToolSourcesCard/>);

    expect(await screen.findByText(/Naming conflict/)).toBeInTheDocument();
    expect(screen.getByText(/list_backups/)).toBeInTheDocument();
  });

  it('renders every registered source, local and external together', async () => {
    mockedGetMCPSources.mockResolvedValue({
      ...baseResponse,
      sources: [
        {name: 'kubectl', kind: 'local', capabilities: ['kubectl'], health: {state: 'healthy'}},
        {name: 'velero-mcp', kind: 'external', capabilities: ['list_backups'], health: {state: 'healthy'}},
      ],
    });

    render(<ToolSourcesCard/>);

    expect(await screen.findAllByText('kubectl')).toHaveLength(2);
    expect(screen.getByText('velero-mcp')).toBeInTheDocument();
  });
});
