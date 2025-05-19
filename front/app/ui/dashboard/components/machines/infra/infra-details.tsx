import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {roboto} from "@/fonts";
import React from "react";
import {
  Text,
  Card,
  GridCol,
  Group,
  SimpleGrid,
  Space,
  Stack,
  Tabs,
  TabsPanel
} from "@mantine/core";
import {IconCheck, IconX} from "@tabler/icons-react";
import AITroubleshooting from "@/app/ui/dashboard/components/machines/infra/ai-troubleshooting";
import Specification from "@/app/ui/dashboard/components/machines/infra/specification";

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
                  machine.status.ready
                    ? <IconCheck size={40} color="teal"/>
                    : <IconX color="red" size={40}/>
                }
                <Text className="text-bold" fw={700}>{machine.metadata?.name}</Text>
              </Group>
            </div>
            <div>
            <Group justify="flex-end">
                <Stack gap="sm" justify="center">
                  <Text size="sm">Namespace</Text>
                  <Text size="xl">
                    {machine.metadata?.namespace}
                  </Text>
                </Stack>
                <Stack gap="sm" justify="center">
                  <Text size="sm">Created</Text>
                  <Text size="xl">
                    {machine.age}
                  </Text>
                </Stack>
            </Group>
            </div>
          </SimpleGrid>
        </Card>
        <Space h="md" />
        <Tabs mb="md" color="#48654a" defaultValue="spec">
          <Tabs.List>
            <Tabs.Tab value="spec">Specification</Tabs.Tab>
            <Tabs.Tab value="troubleshooting">AI Troubleshooting</Tabs.Tab>
          </Tabs.List>
        <TabsPanel value="spec">
          <Space h="lg" />
          <Specification machine={machine} />
        </TabsPanel>
        <TabsPanel value="troubleshooting">
          <Space h="lg" />
          <AITroubleshooting machine={machine} />
        </TabsPanel>
        </Tabs>
      </GridCol>
  )
}
