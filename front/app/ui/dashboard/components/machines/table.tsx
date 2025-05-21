'use client';

import React from "react";
import {Badge, Indicator, Table} from '@mantine/core';
import { GridCol } from '@mantine/core';
import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import {roboto} from "@/fonts";

export default function MachinesTable({
  machines, select
}: {
  machines: MachineType[],
  select: (machine: MachineType) => void
}) {
  return (
    <GridCol span={12}>
      <Table highlightOnHover>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Namespace</Table.Th>
            <Table.Th>Version</Table.Th>
            <Table.Th>Cluster</Table.Th>
            <Table.Th>Age</Table.Th>
            <Table.Th>Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            machines.map( (machine: MachineType, i) => (
              <Table.Tr className={roboto.className} key={i}>
                <Table.Td>
                  <a className="cursor-pointer hover:opacity-70" onClick={() => select(machine)}>{machine.metadata.name}</a>
                </Table.Td>
                <Table.Td>
                  <Badge variant="light" color="gray"> {machine.metadata.namespace} </Badge>
                </Table.Td>
                <Table.Td>{machine.version}</Table.Td>
                <Table.Td>{machine.cluster}</Table.Td>
                <Table.Td ta="center">{machine.age}</Table.Td>
                <Table.Td ta="center">
                  {machine.status?.bootstrapReady && machine.status?.infrastructureReady
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