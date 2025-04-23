'use client';

import React from "react";
import { Table } from '@mantine/core';
import { GridCol } from '@mantine/core';

type MachineDeployment = {
  name: string,
  namespace: string,
  phase: string,
}

export default function MDTable({
  mds,
}: {
  mds: MachineDeployment[]
}) {
  return (
    <GridCol span={12}>
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
    </GridCol>
  );
}