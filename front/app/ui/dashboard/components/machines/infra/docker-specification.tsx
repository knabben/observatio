import {MachineInfraDockerType} from "@/app/ui/dashboard/components/machines/types";
import Panel from "@/app/ui/dashboard/utils/panel";
import {Table} from "@mantine/core";
import React from "react";

export default function DockerSpecification({
  machine,
}: { machine: MachineInfraDockerType }) {
  return (
    <Panel title="Specification" content={
      <Table variant="vertical">
        <Table.Tbody className="text-sm">
          <Table.Tr>
            <Table.Th w={260}>Provider ID</Table.Th>
            <Table.Td>{machine.providerID ?? '—'}</Table.Td>
          </Table.Tr>
          <Table.Tr>
            <Table.Th w={260}>Ready</Table.Th>
            <Table.Td>{machine.ready ? 'true' : 'false'}</Table.Td>
          </Table.Tr>
        </Table.Tbody>
      </Table>
    }/>
  )
}
