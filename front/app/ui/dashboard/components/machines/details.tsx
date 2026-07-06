import React from "react";

import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import ObjectDetails from "@/app/ui/dashboard/base/details";
import Specification from "@/app/ui/dashboard/components/machines/specification";

import {IconCheck, IconX } from "@tabler/icons-react";
import {Group, SimpleGrid, Stack, Text} from "@mantine/core";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";

function buildContext(machine: MachineType): ObjectContext {
  const ready = Boolean(machine.status?.infrastructureReady && machine.status?.bootstrapReady);
  return {
    kind: 'Machine',
    name: machine.metadata?.name ?? '',
    namespace: machine.metadata?.namespace ?? '',
    status: ready ? 'Ready' : 'Not ready',
    keySpecFields: {
      nodeName: machine.nodeName ?? '—',
      providerID: machine.providerID ?? '—',
      version: machine.version ?? '—',
    },
  };
}

export default function MachineDetails({
  machine,
}: { machine: MachineType}) {
  useCurrentObjectContext(buildContext(machine));

  const tabs = [
  {
    label: "Specification",
    content: (machine: MachineType) => <Specification machine={machine} />
  }];
  const headerRender = (machine: MachineType) => (
    <SimpleGrid cols={{base: 1, sm: 2}}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            machine.status?.infrastructureReady && machine.status?.bootstrapReady
            ? <IconCheck size={40} color="teal"/>
            : <IconX color="red" size={40}/>
          }
          <Text className="font-bold" fw={700}>{machine.metadata?.name}</Text>
          <AskAIButton context={buildContext(machine)}/>
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
