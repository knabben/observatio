'use client';

import {SimpleGrid, Table, Text} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import React from "react";
import {MachineType} from "@/app/ui/dashboard/components/machines/types";

export default function Specification({
  machine,
 }: {machine: MachineType}) {
  return (
    <SimpleGrid cols={2}>
      <div>
        <Panel title="Specification" content={
          <Table variant="vertical">
            <Table.Tbody className="text-sm">
              {
                machine.metadata?.ownerReferences.length > 0 &&
                <Table.Tr>
                  <Table.Th w={260}>Owner</Table.Th>
                  <Table.Td>
                    {machine.metadata?.ownerReferences[0].name}
                  </Table.Td>
                </Table.Tr>
              }
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Namespace</Text></Table.Th>
                <Table.Td>
                  <Text style={{ wordBreak: 'break-all' }}>{machine.metadata?.namespace}</Text>
                </Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Cluster</Text></Table.Th>
                <Table.Td>{machine.cluster}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Bootstrap</Text></Table.Th>
                <Table.Td>{machine.bootstrap}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Node Name</Text></Table.Th>
                <Table.Td>{machine.nodeName}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Provider ID</Text></Table.Th>
                <Table.Td>{machine.providerID}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Version</Text></Table.Th>
                <Table.Td>{machine.version}</Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
      </div>
    </SimpleGrid>
  )
}