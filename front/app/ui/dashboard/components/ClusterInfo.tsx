'use client';

import { Table } from '@mantine/core';

type Service = {
  name: string,
  path: string
}

// ClusterInfo: Cluster details and access URLs.
export default function ClusterInfo({
  clusterInfo,
}: {
  clusterInfo: Service[]
}) {
  return (
    <Table striped highlightOnHover withTableBorder withColumnBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>URL</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {
          clusterInfo.map( (service) => (
            <Table.Tr key={service.name}>
              <Table.Td>{service.name}</Table.Td>
              <Table.Td><a href={service.path}>{service.path}</a></Table.Td>
            </Table.Tr>
          ))
        }
      </Table.Tbody>
    </Table>
  );
}