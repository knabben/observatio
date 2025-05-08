// details.test.tsx
import React from 'react';
import {render, screen} from '@testing-library/react';
import MachineDetails, {MachineType} from './details';

describe('MachineDetails Component', () => {
  const mockMachine: MachineType = {
    name: 'TestMachine',
    namespace: 'default',
    owner: 'admin',
    bootstrap: 'init-bootstrap',
    cluster: 'test-cluster',
    nodeName: 'node1',
    providerID: 'provider-1234',
    version: 'v1.23.0',
    created: '2023-10-01',
    bootstrapReady: true,
    infrastructureReady: true,
    phase: 'Running',
  };

  it('renders machine details with correct data', () => {
    render(<MachineDetails machine={mockMachine}/>);

    expect(screen.getByText(/Name:/)).toBeInTheDocument();
    expect(screen.getByText(/TestMachine/i)).toBeInTheDocument();
    expect(screen.getByText(/Phase:/)).toBeInTheDocument();
    expect(screen.getByText(/Running/i)).toBeInTheDocument();
    expect(screen.getByText(/Age:/)).toBeInTheDocument();
    expect(screen.getByText(/2023-10-01/i)).toBeInTheDocument();
  });

  it('renders green indicator when both infrastructureReady and bootstrapReady are true', () => {
    render(<MachineDetails machine={mockMachine}/>);

    const indicator = screen.getByText(/TestMachine/i).closest('div');
    expect(indicator).toHaveClass('mantine-Indicator');
    expect(indicator?.querySelector('.mantine-Indicator-indicator')).toHaveStyle('background-color: green');
  });

  it('renders red indicator when infrastructureReady or bootstrapReady is false', () => {
    const updatedMachine = {...mockMachine, bootstrapReady: false};
    render(<MachineDetails machine={updatedMachine}/>);

    const indicator = screen.getByText(/TestMachine/i).closest('div');
    expect(indicator).toHaveClass('mantine-Indicator');
    expect(indicator?.querySelector('.mantine-Indicator-indicator')).toHaveStyle('background-color: red');
  });

  it('renders the machine specification table correctly', () => {
    render(<MachineDetails machine={mockMachine}/>);

    expect(screen.getByText(/Namespace/i)).toBeInTheDocument();
    expect(screen.getByText(/default/i)).toBeInTheDocument();
    expect(screen.getByText(/Cluster/i)).toBeInTheDocument();
    expect(screen.getByText(/test-cluster/i)).toBeInTheDocument();
    expect(screen.getByText(/Owner/i)).toBeInTheDocument();
    expect(screen.getByText(/admin/i)).toBeInTheDocument();
    expect(screen.getByText(/Bootstrap/i)).toBeInTheDocument();
    expect(screen.getByText(/init-bootstrap/i)).toBeInTheDocument();
    expect(screen.getByText(/Node/i)).toBeInTheDocument();
    expect(screen.getByText(/node1/i)).toBeInTheDocument();
    expect(screen.getByText(/ProviderID/i)).toBeInTheDocument();
    expect(screen.getByText(/provider-1234/i)).toBeInTheDocument();
    expect(screen.getByText(/Version/i)).toBeInTheDocument();
    expect(screen.getByText(/v1.23.0/i)).toBeInTheDocument();
  });
});