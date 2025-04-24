'use client';

import { Table } from '@mantine/core';

type Component = {
    name: string,
    kind: string,
    version: string,
}

// Versions: Cluster components versions and enumeration.
export default function Versions({
    components,
}: {
    components: Component[]
}) {
  return (
  <Table striped highlightOnHover withTableBorder withColumnBorders>
    <Table.Thead>
    <Table.Tr>
      <Table.Th>Name</Table.Th>
      <Table.Th>Kind</Table.Th>
      <Table.Th>Versions</Table.Th>
    </Table.Tr>
    </Table.Thead>
    <Table.Tbody>
      {
        components.map((component) => (
          <Table.Tr key={component.name}>
            <Table.Td>{component.name}</Table.Td>
            <Table.Td>{component.kind}</Table.Td>
            <Table.Td>{component.version}</Table.Td>
          </Table.Tr>
        ))
      }
    </Table.Tbody>
  </Table>
  );
}
