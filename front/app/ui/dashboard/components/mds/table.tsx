'use client';

import React from "react";
import { Table, Indicator } from '@mantine/core';
import { GridCol } from '@mantine/core';

import {MachineDeploymentType} from '@/app/ui/dashboard/components/mds/types';

export default function MDTable({
  mds,
}: {
  mds: MachineDeploymentType[]
}) {
  return (
    <GridCol span={12}>
      <Table highlightOnHover>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Namespace</Table.Th>
            <Table.Th>Replicas</Table.Th>
            <Table.Th>Cluster</Table.Th>
            <Table.Th>Age</Table.Th>
            <Table.Th ta="center">Phase</Table.Th>
            <Table.Th ta="center">Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            mds.map((md, i) => (
              <Table.Tr key={i}>
                <Table.Td>{md.name}</Table.Td>
                <Table.Td>{md.namespace}</Table.Td>
                <Table.Td>{md.replicas}</Table.Td>
                <Table.Td>{md.cluster}</Table.Td>
                <Table.Td>{md.created}</Table.Td>
                <Table.Td ta="center">{md.phase}</Table.Td>
                <Table.Td ta="center">
                  {md.unavailableReplicas == 0
                    ? <Indicator inline processing color="green" size={15}/>
                    : <Indicator inline processing color="red" size={15}/>
                  }
                </Table.Td>
              </Table.Tr>
            ))
          }
        </Table.Tbody>
      </Table>
    </GridCol>
  );
}