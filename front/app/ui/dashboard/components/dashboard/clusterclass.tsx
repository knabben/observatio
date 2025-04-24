'use client';

import { Table } from '@mantine/core';

type Conditions = {
  type: string,
  status: boolean,
  lastTransitionTime: string,
}

type Clusterclass = {
  name: string,
  namespace: string,
  generation: bigint,
  conditions: Conditions[]
}

// ClusterClass: details the cluster classes existent in the cluster.
export default function ClusterClass({
  clusterClass,
}: {
  clusterClass: Clusterclass[]
}) {
  return (
    <Table striped highlightOnHover withTableBorder withColumnBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>Namespace</Table.Th>
          <Table.Th>Updates</Table.Th>
          <Table.Th>Status</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {
          clusterClass.map( (cc, i) => (
            <Table.Tr key={i}>
              <Table.Td>{cc.name}</Table.Td>
              <Table.Td>{cc.namespace}</Table.Td>
              <Table.Td>{cc.generation}</Table.Td>
              {cc.conditions.map((condition, i) => (
                  <Table.Td key={i} rowSpan={1}>
                    {condition.type} - {condition.status}
                  </Table.Td>
              ))}
            </Table.Tr>
          ))
        }
      </Table.Tbody>
    </Table>
  );
}