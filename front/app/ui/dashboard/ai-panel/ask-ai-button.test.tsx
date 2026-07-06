import '@testing-library/jest-dom';
import {fireEvent, screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {AIPanelProvider, useAIPanel} from './ai-panel-context';
import {AskAIButton} from './ask-ai-button';

function Probe() {
  const {isOpen, queryField} = useAIPanel();
  return <div data-testid="probe">{isOpen ? 'open' : 'closed'}::{queryField}</div>;
}

describe('AskAIButton', () => {
  it('opens the panel pre-filled with this screen\'s object context in one click', () => {
    render(
      <AIPanelProvider>
        <AskAIButton context={{
          kind: 'Machine', name: 'm1', namespace: 'default', status: 'Ready', keySpecFields: {version: 'v1.30.0'},
        }}/>
        <Probe/>
      </AIPanelProvider>,
    );

    expect(screen.getByTestId('probe')).toHaveTextContent('closed::');

    fireEvent.click(screen.getByRole('button', {name: /ask ai about this/i}));

    const probe = screen.getByTestId('probe');
    expect(probe).toHaveTextContent(/^open::/);
    expect(probe).toHaveTextContent('Machine "m1"');
    expect(probe).toHaveTextContent('version=v1.30.0');
  });
});
