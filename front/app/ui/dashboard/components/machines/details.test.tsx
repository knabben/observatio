// details.test.tsx
import React from 'react';
import {screen} from '@testing-library/react';
import {render} from "@/app/ui/dashboard/utils/test-render";
import MachineDetails from './details';
import {MachineType} from './types'
import {describe, expect, it, jest} from "@jest/globals";

jest.mock('next/font/google', () => ({
  Source_Sans_3: () => ({
    style: {
      fontFamily: 'Source Sans 3',
    },
  }),
  Open_Sans: () => ({
    style: {
      fontFamily: 'Open Sans',
    },
  }),
  Lora: () => ({
    style: {
      fontFamily: 'Lora',
    },
  }),
  Roboto: () => ({
    style: {
      fontFamily: 'Roboto',
    },
  }),
  Inter: () => ({
    style: {
      fontFamily: 'Inter',
    },
  }),
}));

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

    expect(screen.getByText(/Name:/)).toBe("TestMachine")
  });
});