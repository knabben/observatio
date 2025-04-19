'use client';

import { Table } from '@mantine/core';

type Machine = {
  name: string,
  namespace: string,
  version: string,
  nodeName: string,
  phase: string,
}

// MachineLister: details the machines existent in the cluster.
export default function MachineLister ({
  machines,
}: {
  machines: Machine[]
}) {
  return (
    <Table striped highlightOnHover withTableBorder withColumnBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>Namespace</Table.Th>
          <Table.Th>Version</Table.Th>
          <Table.Th>Node</Table.Th>
          <Table.Th>Phase</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {
          machines.map( (machine, i) => (
            <Table.Tr key={i}>
              <Table.Td>{machine.name}</Table.Td>
              <Table.Td>{machine.namespace}</Table.Td>
              <Table.Td>{machine.version}</Table.Td>
              <Table.Td>{machine.nodeName}</Table.Td>
              <Table.Td>{machine.phase}</Table.Td>
            </Table.Tr>
          ))
        }
      </Table.Tbody>
    </Table>
  );
}