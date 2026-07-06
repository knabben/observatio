import React from "react";
import {MachineInfraDockerType} from "@/app/ui/dashboard/components/machines/types";
import {
  Text,
  Group,
  SimpleGrid,
  Stack,
} from "@mantine/core";
import {IconCheck, IconX} from "@tabler/icons-react";
import AITroubleshooting from "@/app/ui/dashboard/base/ai-troubleshooting";
import ObjectDetails from "@/app/ui/dashboard/base/details";

export default function MachineInfraDockerDetails({
  machine
}: {machine: MachineInfraDockerType}) {
  const tabs = [
    {
      label: "AI Troubleshooting",
      content: (machine: MachineInfraDockerType) => <AITroubleshooting
        objectType="dockermachine"
        objectName={machine.metadata?.name ?? ''}
        objectNamespace={machine.metadata?.namespace ?? ''}
        conditions={[]} />
    }];

  const headerRender = (machine: MachineInfraDockerType) => (
    <SimpleGrid cols={{base: 1, sm: 2}}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            machine.ready
              ? <IconCheck size={40} color="teal"/>
              : <IconX color="red" size={40}/>
          }
          <Text className="font-bold" fw={700}>{machine.metadata?.name}</Text>
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
            <Text size="sm">Age</Text>
            <Text size="xl">
              {machine.age ?? '—'}
            </Text>
          </Stack>
        </Group>
      </div>
    </SimpleGrid>
  );

  return (
    <ObjectDetails
      object={machine}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
