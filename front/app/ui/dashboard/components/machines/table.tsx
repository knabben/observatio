'use client';

import React from "react";
import { Indicator, Table } from '@mantine/core';
import { GridCol } from '@mantine/core';
import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import {roboto} from "@/fonts";

export default function MachinesTable({
  machines,
  loading,
}: {
  machines: MachineType[]
  loading: boolean
}) {
  return (
    <GridCol span={12}>
      <Table highlightOnHover>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Namespace</Table.Th>
            <Table.Th>Version</Table.Th>
            <Table.Th>Node</Table.Th>
            <Table.Th>Cluster</Table.Th>
            <Table.Th ta="center">Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            machines.map( (machine, i) => (
              <Table.Tr className={roboto.className} key={i}>
                <Table.Td>{machine.name}</Table.Td>
                <Table.Td>{machine.namespace}</Table.Td>
                <Table.Td>{machine.version}</Table.Td>
                <Table.Td>{machine.nodeName}</Table.Td>
                <Table.Td>{machine.cluster}</Table.Td>
                <Table.Td ta="center">
                  {machine.bootstrapReady && machine.infrastructureReady
                    ? <Indicator inline processing color="green" size={15}/>
                    : <Indicator inline processing color="red" size={15}/>
                  }</Table.Td>
              </Table.Tr>
            ))
          }
        </Table.Tbody>
      </Table>
    </GridCol>
  )
}