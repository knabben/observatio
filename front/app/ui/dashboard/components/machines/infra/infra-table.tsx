import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {roboto} from "@/fonts";
import React from "react";
import {GridCol, Indicator, Table, Badge} from "@mantine/core";


export default function MachineInfraTable({
  machines, select
}: {
  machines: MachineInfraType[],
  select: (machine: MachineInfraType) => void
}) {
  return (
    <GridCol span={12}>
      <Table highlightOnHover>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Namespace</Table.Th>
            <Table.Th>ProviderID</Table.Th>
            <Table.Th>Template</Table.Th>
            <Table.Th>Age</Table.Th>
            <Table.Th>Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            machines?.map( (machine: MachineInfraType) => (
              <Table.Tr className={roboto.className} key={machine.metadata.name}>
                <Table.Td>
                  <a className="cursor-pointer hover:opacity-70" onClick={() => select(machine)}>{machine.metadata.name}</a>
                </Table.Td>
                <Table.Td>
                  <Badge variant="light" color="gray"> {machine.metadata.namespace} </Badge>
                </Table.Td>
                <Table.Td>{machine.providerID}</Table.Td>
                <Table.Td>{machine.template}</Table.Td>
                <Table.Td>{machine.age}</Table.Td>
                <Table.Td className="text-center align-middle">
                {
                  machine.status.ready
                  ? <Indicator inline processing color="green" size={22}/>
                  : <Indicator inline processing color="red" size={22}/>
                }
                </Table.Td>
              </Table.Tr>
            ))
          }
        </Table.Tbody>
      </Table>
    </GridCol>
  )
}