import React from "react";
import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {
  Text,
  Group,
  SimpleGrid,
  Stack,
} from "@mantine/core";
import {IconCheck, IconX} from "@tabler/icons-react";
import Specification from "@/app/ui/dashboard/components/machines/infra/specification";
import ObjectDetails from "@/app/ui/dashboard/base/details";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function buildContext(machine: MachineInfraType): ObjectContext {
  return {
    kind: 'VSphereMachine',
    name: machine.metadata?.name ?? '',
    namespace: machine.metadata?.namespace ?? '',
    status: machine.status?.ready ? 'Ready' : `Not ready${machine.status?.failureReason ? `: ${machine.status.failureReason}` : ''}`,
    keySpecFields: {
      template: machine.template ?? '—',
      numCPUs: String(machine.numCPUs ?? '—'),
      memoryMiB: String(machine.memoryMiB ?? '—'),
    },
  };
}

export default function MachineInfraDetails({
  machine
}: {machine: MachineInfraType}) {
  useCurrentObjectContext(buildContext(machine));

  const tabs = [
    {
      label: "Specification",
      content: (machine: MachineInfraType) => <Specification machine={machine} />
    },
    {
      label: "YAML",
      content: (machine: MachineInfraType) => <ObjectTree
        gvr={RESOURCE_GVR.vsphereMachine}
        namespace={machine.metadata?.namespace ?? ''}
        name={machine.metadata?.name ?? ''}
        resourceVersion={machine.metadata?.resourceVersion}
      />
    }];

  const headerRender = (machine: MachineInfraType) => (
    <SimpleGrid cols={{base: 1, sm: 2}}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            machine.status?.ready
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
