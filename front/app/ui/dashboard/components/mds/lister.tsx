'use client';

import { Table } from '@mantine/core';

type MachineDeploymentLister = {
  name: string,
  namespace: string,
  phase: string,
}

// MachineDeploymentLister: details the mds existent in the cluster.
export default function MachineDeploymentLister({
  mds,
}: {
  mds: MachineDeploymentLister[]
}) {
  return (
    <Table striped highlightOnHover withTableBorder withColumnBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>Namespace</Table.Th>
          <Table.Th>Phase</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {
          mds.map( (md, i) => (
            <Table.Tr key={i}>
              <Table.Td>{md.name}</Table.Td>
              <Table.Td>{md.namespace}</Table.Td>
              <Table.Td>{md.phase}</Table.Td>
            </Table.Tr>
          ))
        }
      </Table.Tbody>
    </Table>
  );
}