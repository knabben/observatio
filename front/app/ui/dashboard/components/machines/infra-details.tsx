import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {roboto} from "@/fonts";
import Panel from "@/app/ui/dashboard/utils/panel";
import React from "react";
import {
  Text,
  Card,
  Chip,
  GridCol,
  Group,
  SimpleGrid,
  Space,
  Stack,
  Table,
  Tabs,
  TabsPanel
} from "@mantine/core";
import { XMarkIcon } from '@heroicons/react/24/outline';
import {IconCheck, IconCpu, IconDatabase, IconDeviceFloppy, IconX} from "@tabler/icons-react";

export default function MachineInfraDetails({
  machine
}: {machine: MachineInfraType}) {
  return (
      <GridCol className={roboto.className} span={12}>
        <Card withBorder shadow="sm" padding="md" radius="md">
          <SimpleGrid cols={2}>
            <div className="flex items-center h-full">
              <Group justify="flex-start">
                {
                  machine.ready
                    ? <IconCheck size={40} color="teal"/>
                    : <IconX color="red" size={40}/>
                }
                <Text className="text-bold" fw={700}>{machine.name}</Text>
              </Group>
            </div>
            <div>
            <Group justify="flex-end">
                <Stack gap="sm" justify="center">
                  <Text size="sm">Namespace</Text>
                  <Text size="xl">
                    {machine.namespace}
                  </Text>
                </Stack>
                <Stack gap="sm" justify="center">
                  <Text size="sm">Created</Text>
                  <Text size="xl">
                    {machine.created}
                  </Text>
                </Stack>
              </Group>
            </div>
          </SimpleGrid>
        </Card>
        <Space h="md" />
        <Tabs mb="md" defaultValue="spec">
          <Tabs.List>
            <Tabs.Tab value="spec">Specification</Tabs.Tab>
            <Tabs.Tab value="troubleshooting">AI Troubleshooting</Tabs.Tab>
          </Tabs.List>
        <TabsPanel value="spec">
          <Space h="lg" />
          <SimpleGrid cols={2}>
          <div>
            <Panel title="Specification" content={
              <Table variant="vertical">
                <Table.Tbody className="text-sm">
                  <Table.Tr>
                    <Table.Th w={260}><Text fw={500}>Namespace</Text></Table.Th>
                    <Table.Td>
                      <Text style={{ wordBreak: 'break-all' }}>{machine.namespace}</Text>
                    </Table.Td>
                  </Table.Tr>
                  <Table.Tr>
                    <Table.Th w={260}><Text fw={500}>Provider</Text></Table.Th>
                    <Table.Td>{machine.providerID}</Table.Td>
                  </Table.Tr>
                  { machine.failureDomain &&
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
                    { machine.numCoresPerSocket &&
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
                    <Table.Tr>
                      <Table.Th w={260}><Text fw={500}>Failure Reason</Text></Table.Th>
                      <Table.Td>{machine.failureReason}</Table.Td>
                    </Table.Tr>
                    <Table.Tr>
                      <Table.Th w={260}><Text fw={500}>Failure Message</Text></Table.Th>
                      <Table.Td>{machine.failureMessage}</Table.Td>
                    </Table.Tr>
                  </Table.Tbody>
                </Table>
              } />
            </div>
          </SimpleGrid>
        </TabsPanel>
        <TabsPanel value="troubleshooting">
          <Space h="lg" />
          <Panel title="Machine conditions" content={
            <Table variant="vertical">
              <Table.Tbody className="text-sm">
                {
                  machine.conditions?.map((condition, ic) => (
                    <Table.Tr key={ic}>
                      <Table.Td>
                        {
                          condition.status.toLowerCase() == "true"
                            ? <Chip key={ic} className="p-1" defaultChecked color="teal" variant="light">{condition.type}</Chip>
                            : <Chip key={ic} defaultChecked icon={<XMarkIcon />} color="red" variant="light">{condition.type}</Chip>
                        }
                      </Table.Td>
                      <Table.Td>{condition.lastTransitionTime}</Table.Td>
                    </Table.Tr>
                  ))
                }
              </Table.Tbody>
            </Table>
          } />
        </TabsPanel>
        </Tabs>
      </GridCol>
  )
}
