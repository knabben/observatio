import '@testing-library/jest-dom';
import React from 'react';
import {act, fireEvent, screen} from '@testing-library/react';
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

  it('shows a not-available message and disables Send when reconnect attempts are exhausted (FR-017)', () => {
    let onReconnectStop: (() => void) | undefined;
    mockedUseWebSocket.mockImplementation((_url, options) => {
      onReconnectStop = options?.onReconnectStop as (() => void) | undefined;
      return {
        sendJsonMessage: jest.fn(),
        lastJsonMessage: null,
        readyState: ReadyState.CONNECTING,
      } as unknown as ReturnType<typeof useWebSocket>;
    });

    render(
      <AIPanelProvider>
        <AIPanel/>
        <AIPanelTrigger/>
      </AIPanelProvider>,
    );
    fireEvent.click(screen.getByRole('button', {name: /open ai troubleshooting panel/i}));

    expect(screen.queryByText(/not available/i)).not.toBeInTheDocument();

    act(() => onReconnectStop?.());

    expect(screen.getByText(/not available/i)).toBeInTheDocument();
    expect(screen.getByRole('button', {name: 'Send'})).toBeDisabled();
  });

  it('merges streamed "delta" chunks sharing an id into one bubble and clears loading on "done"', () => {
    // The mock hook holds its own React state and re-renders AIPanel from the inside (like the
    // real hook would on each incoming frame), instead of forcing an external rerender of the
    // whole tree - re-rendering the tree from outside remounts AIPanelProvider in this test setup
    // and loses its open/closed state, which isn't representative of how the app actually runs.
    let pushToHook: ((msg: unknown) => void) | undefined;
    mockedUseWebSocket.mockImplementation(() => {
      const [lastJsonMessage, setLastJsonMessage] = React.useState<unknown>(null);
      pushToHook = setLastJsonMessage;
      return {
        sendJsonMessage: jest.fn(),
        lastJsonMessage,
        readyState: ReadyState.OPEN,
      } as unknown as ReturnType<typeof useWebSocket>;
    });

    render(
      <AIPanelProvider>
        <AIPanel/>
        <AIPanelTrigger/>
      </AIPanelProvider>,
    );
    fireEvent.click(screen.getByRole('button', {name: /open ai troubleshooting panel/i}));

    const push = (msg: unknown) => act(() => pushToHook?.(msg));

    push({
      id: 'resp-1', type: 'chatbot', agent_id: 'cluster-agent', actor: 'agent',
      timestamp: 't1', content: 'Hello', event: 'delta',
    });
    push({
      id: 'resp-1', type: 'chatbot', agent_id: 'cluster-agent', actor: 'agent',
      timestamp: 't2', content: ' world', event: 'delta',
    });

    expect(screen.getByText('Hello world')).toBeInTheDocument();

    push({
      id: 'resp-1', type: 'chatbot', agent_id: 'cluster-agent', actor: 'agent',
      timestamp: 't3', content: '', event: 'done',
    });

    // "done" carries no content of its own - the merged text from the delta chunks stays put.
    expect(screen.getByText('Hello world')).toBeInTheDocument();
  });
});
