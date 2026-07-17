import '@testing-library/jest-dom';
import {screen} from '@testing-library/react';
import {render} from '@/app/ui/dashboard/utils/test-render';
import {SeverityBanner} from './severity-banner';
import {FailureSeverity} from '@/app/ui/dashboard/shared/use-day2-ops';

const selfHealing: FailureSeverity = {objectRef: null, level: 'self_healing', reason: 'MachineHealthCheck is remediating 1 unhealthy machine(s)'};
const needsInvestigation: FailureSeverity = {objectRef: null, level: 'needs_investigation', reason: 'maxUnhealthy threshold breached'};
const managementCritical: FailureSeverity = {objectRef: null, level: 'management_critical', reason: 'API server unreachable'};

describe('SeverityBanner', () => {
  it('renders nothing when there are no severities', () => {
    render(<SeverityBanner severities={[]}/>);
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('shows self-healing informationally, not as an alert-level banner', () => {
    render(<SeverityBanner severities={[selfHealing]}/>);
    expect(screen.getByText(/remediating/i)).toBeInTheDocument();
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('shows needs-investigation as a warning banner', () => {
    render(<SeverityBanner severities={[needsInvestigation]}/>);
    expect(screen.getByText(/maxUnhealthy threshold breached/i)).toBeInTheDocument();
  });

  it('shows the highest-severity banner when multiple are present', () => {
    render(<SeverityBanner severities={[selfHealing, needsInvestigation, managementCritical]}/>);
    expect(screen.getByText(/API server unreachable/i)).toBeInTheDocument();
    expect(screen.queryByText(/maxUnhealthy threshold breached/i)).not.toBeInTheDocument();
  });

  it('never shows a self-healing event with the same styling as an actionable failure', () => {
    const {container: selfHealingContainer} = render(<SeverityBanner severities={[selfHealing]}/>);
    const {container: criticalContainer} = render(<SeverityBanner severities={[managementCritical]}/>);
    expect(selfHealingContainer.querySelector('[role="alert"]')).toBeNull();
    expect(criticalContainer.querySelector('[role="alert"]')).not.toBeNull();
  });

  it('renders the enriched recoverable reason from a CA-secret-missing severity (008/US2)', () => {
    const recoverable: FailureSeverity = {
      objectRef: null, level: 'management_critical',
      reason: "Cluster prod-1's CA secret is missing or inaccessible. A backup completed 3h0m0s ago covers this cluster — recovery is straightforward.",
      recoveryInfo: {recoverable: true, coveringBackupAge: '3h0m0s'},
    };
    render(<SeverityBanner severities={[recoverable]}/>);
    expect(screen.getByText(/recovery is straightforward/i)).toBeInTheDocument();
  });

  it('renders the enriched unrecoverable reason from a CA-secret-missing severity (008/US2)', () => {
    const unrecoverable: FailureSeverity = {
      objectRef: null, level: 'management_critical',
      reason: "Cluster prod-1's CA secret is missing or inaccessible. No covering backup exists — this cluster's data is unrecoverable.",
      recoveryInfo: {recoverable: false},
    };
    render(<SeverityBanner severities={[unrecoverable]}/>);
    expect(screen.getByText(/unrecoverable/i)).toBeInTheDocument();
  });
});
