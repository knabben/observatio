'use client';

import {Group, SimpleGrid, Table, Text} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import {IconCpu, IconDatabase, IconDeviceFloppy} from "@tabler/icons-react";
import React from "react";
import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";

export default function Specification({
  machine,
 }: {machine: MachineInfraType}) {
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
                <Table.Th w={260}><Text fw={500}>Provider</Text></Table.Th>
                <Table.Td>{machine.providerID}</Table.Td>
              </Table.Tr>
              {
                machine.failureDomain &&
                <Table.Tr>
                  <Table.Th w={260}><Text fw={500}>Failure Domain</Text></Table.Th>
                  <Table.Td>{machine.failureDomain}</Table.Td>
                </Table.Tr>
              }
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Power Off Mode</Text></Table.Th>
                <Table.Td>{machine.powerOffMode}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Template</Text></Table.Th>
                <Table.Td>{machine.template}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}>Clone Mode</Text></Table.Th>
                <Table.Td>{machine.cloneMode}</Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
      </div>
      <div>
        <Panel title="Machine details" content={
          <Table variant="vertical">
            <Table.Tbody className="text-sm">
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}><Group><IconCpu />CPUs</Group></Text></Table.Th>
                <Table.Td>{machine.numCPUs}</Table.Td>
              </Table.Tr>
              {
                machine.numCoresPerSocket &&
                <Table.Tr>
                  <Table.Th>CPU Per Socket</Table.Th>
                  <Table.Td>{machine.numCoresPerSocket}</Table.Td>
                </Table.Tr>
              }
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}><Group><IconDatabase />Memory</Group></Text></Table.Th>
                <Table.Td>{machine.memoryMiB} MiB</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}><Text fw={500}><Group><IconDeviceFloppy />Disk size</Group></Text></Table.Th>
                <Table.Td>{machine.diskGiB} GiB</Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
      </div>
    </SimpleGrid>
  )
}