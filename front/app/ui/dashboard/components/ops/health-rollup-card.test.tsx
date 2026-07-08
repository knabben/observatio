import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {HealthRollupCard} from './health-rollup-card';
import {HealthRollup} from '@/app/ui/dashboard/shared/use-day2-ops';

const baseRollup: HealthRollup = {
  category: 'cluster',
  healthy: 0,
  degraded: 0,
  failed: 0,
  unavailable: false,
};

describe('HealthRollupCard', () => {
  it('renders healthy/degraded/failed counts', () => {
    render(<HealthRollupCard rollup={{...baseRollup, healthy: 3, degraded: 1, failed: 2}}/>);

    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
  });

  it('shows an "all clear" state when nothing is degraded or failed', () => {
    render(<HealthRollupCard rollup={{...baseRollup, healthy: 5}}/>);

    expect(screen.getByText(/all clear/i)).toBeInTheDocument();
  });

  it('shows a "data unavailable" state when the category is unavailable', () => {
    render(<HealthRollupCard rollup={{...baseRollup, unavailable: true}}/>);

    expect(screen.getByText(/data unavailable/i)).toBeInTheDocument();
    expect(screen.queryByText(/all clear/i)).not.toBeInTheDocument();
  });

  it('renders the human-readable category title', () => {
    render(<HealthRollupCard rollup={{...baseRollup, category: 'machine_deployment', healthy: 1}}/>);

    expect(screen.getByText('Machine Deployments')).toBeInTheDocument();
  });
});
