import React from "react";

import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import ObjectDetails from "@/app/ui/dashboard/base/details";
import Specification from "@/app/ui/dashboard/components/machines/specification";
import AITroubleshooting from "@/app/ui/dashboard/base/ai-troubleshooting";

import {IconCheck, IconX } from "@tabler/icons-react";
import {Group, SimpleGrid, Stack, Text} from "@mantine/core";

export default function MachineDetails({
  machine,
}: { machine: MachineType}) {
  const tabs = [
  {
    label: "Specification",
    content: (machine: MachineType) => <Specification machine={machine} />
  },
  {
    label: "AI Troubleshooting",
    content: (machine: MachineType) => <AITroubleshooting
      objectType="machine"
      objectName={machine.metadata.name}
      objectNamespace={machine.metadata.namespace}
      conditions={machine.status.conditions} />
  }];
  const headerRender = (machine: MachineType) => (
    <SimpleGrid cols={2}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            machine.status.infrastructureReady && machine.status.bootstrapReady
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
  );

  return (
    <ObjectDetails
      object={machine}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
