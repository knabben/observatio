import '@testing-library/jest-dom';
import {fireEvent, screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import useWebSocket, {ReadyState} from 'react-use-websocket';

jest.mock('react-use-websocket');

import {AIPanelProvider} from './ai-panel-context';
import AIPanel from './ai-panel';
import AIPanelTrigger from './ai-panel-trigger';

const mockedUseWebSocket = useWebSocket as jest.MockedFunction<typeof useWebSocket>;

beforeAll(() => {
  (global as unknown as {ResizeObserver: unknown}).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
});

afterEach(() => jest.resetAllMocks());

function renderOpenPanel(lastJsonMessage: unknown) {
  mockedUseWebSocket.mockReturnValue({
    sendJsonMessage: jest.fn(),
    lastJsonMessage,
    readyState: ReadyState.OPEN,
  } as unknown as ReturnType<typeof useWebSocket>);

  const result = render(
    <AIPanelProvider>
      <AIPanel/>
      <AIPanelTrigger/>
    </AIPanelProvider>,
  );
  fireEvent.click(screen.getByRole('button', {name: /open ai troubleshooting panel/i}));
  return result;
}

describe('AIPanel', () => {
  it('renders AI/user message content as safe plain text, never parsed as HTML', () => {
    const {container} = renderOpenPanel({
      id: '1',
      type: 'chatbot',
      agent_id: 'cluster-agent',
      actor: 'agent',
      timestamp: 'now',
      content: '<img src=x onerror="window.__pwned=true">',
    });
    expect(screen.getByText('<img src=x onerror="window.__pwned=true">')).toBeInTheDocument();
    expect(container.querySelector('img')).toBeNull();
  });
});
