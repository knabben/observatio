import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {roboto} from "@/fonts";
import React from "react";
import {GridCol, Indicator, Table} from "@mantine/core";


export default function MachineInfraTable({
  machines
}: {
  machines: MachineInfraType[]
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
            <Table.Th ta="center">Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            machines?.map( (machine: MachineInfraType) => (
              <Table.Tr className={roboto.className} key={machine.name}>
                <Table.Td>{machine.name}</Table.Td>
                <Table.Td>{machine.namespace}</Table.Td>
                <Table.Td>{machine.providerID}</Table.Td>
                <Table.Td>{machine.template}</Table.Td>
                <Table.Td>{machine.created}</Table.Td>
                <Table.Td ta="center">
                  {machine.ready
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