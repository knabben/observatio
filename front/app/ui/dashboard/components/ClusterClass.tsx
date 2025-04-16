'use client';

import { Table } from '@mantine/core';

type Conditions = {
  type: string,
  status: boolean,
  lastTransitionTime: string,
}

type ClusterClass = {
  name: string,
  namespace: string,
  generation: bigint,
  conditions: Conditions[]
}

// ClusterClass: details the cluster classes existent in the cluster.
export default function ClusterClass({
  clusterClass,
}: {
  clusterClass: ClusterClass[]
}) {
  return (
    <Table striped highlightOnHover withTableBorder withColumnBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>Namespace</Table.Th>
          <Table.Th>Updates</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {
          clusterClass.map( (cc) => (
            <Table.Tr key={cc.name}>
              <Table.Td>{cc.name}</Table.Td>
              <Table.Td>{cc.namespace}</Table.Td>
              <Table.Td>{cc.generation}</Table.Td>
            </Table.Tr>
          ))
        }
      </Table.Tbody>
    </Table>
  );
}