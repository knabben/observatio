import '@testing-library/jest-dom';
import {screen, fireEvent} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import useWebSocket, {ReadyState} from 'react-use-websocket';

jest.mock('react-use-websocket');

import {AIPanelProvider} from './ai-panel-context';
import AIPanel from './ai-panel';
import AIPanelTrigger from './ai-panel-trigger';

const mockedUseWebSocket = useWebSocket as jest.MockedFunction<typeof useWebSocket>;

// Mantine's ScrollArea (used inside AIPanel) measures via ResizeObserver, which jsdom does not implement.
beforeAll(() => {
  (global as unknown as {ResizeObserver: unknown}).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
});

beforeEach(() => {
  mockedUseWebSocket.mockReturnValue({
    sendJsonMessage: jest.fn(),
    lastJsonMessage: null,
    readyState: ReadyState.OPEN,
  } as unknown as ReturnType<typeof useWebSocket>);
});

afterEach(() => jest.resetAllMocks());

function renderShell() {
  return render(
    <AIPanelProvider>
      <div>some screen content</div>
      <AIPanel/>
      <AIPanelTrigger/>
    </AIPanelProvider>,
  );
}

describe('AIPanelTrigger', () => {
  it('opens the same global panel instance when activated, regardless of the surrounding screen', () => {
    renderShell();
    expect(screen.queryByText('AI Troubleshooting')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', {name: /open ai troubleshooting panel/i}));

    expect(screen.getByText('AI Troubleshooting')).toBeInTheDocument();
  });

  it('hides itself while the panel is already open (single entry point, no duplicate trigger)', () => {
    renderShell();
    fireEvent.click(screen.getByRole('button', {name: /open ai troubleshooting panel/i}));
    expect(screen.queryByRole('button', {name: /open ai troubleshooting panel/i})).not.toBeInTheDocument();
  });
});
