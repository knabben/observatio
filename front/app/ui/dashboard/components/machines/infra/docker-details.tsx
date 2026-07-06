import React from "react";
import {MachineInfraDockerType} from "@/app/ui/dashboard/components/machines/types";
import {
  Text,
  Group,
  SimpleGrid,
  Stack,
} from "@mantine/core";
import {IconCheck, IconX} from "@tabler/icons-react";
import DockerSpecification from "@/app/ui/dashboard/components/machines/infra/docker-specification";
import ObjectDetails from "@/app/ui/dashboard/base/details";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function buildContext(machine: MachineInfraDockerType): ObjectContext {
  return {
    kind: 'DockerMachine',
    name: machine.metadata?.name ?? '',
    namespace: machine.metadata?.namespace ?? '',
    status: machine.ready ? 'Ready' : 'Not ready',
    keySpecFields: {
      providerID: machine.providerID ?? '—',
    },
  };
}

export default function MachineInfraDockerDetails({
  machine
}: {machine: MachineInfraDockerType}) {
  useCurrentObjectContext(buildContext(machine));

  const tabs = [
    {
      label: "Specification",
      content: (machine: MachineInfraDockerType) => <DockerSpecification machine={machine} />
    },
    {
      label: "YAML",
      content: (machine: MachineInfraDockerType) => <ObjectTree
        gvr={RESOURCE_GVR.dockerMachine}
        namespace={machine.metadata?.namespace ?? ''}
        name={machine.metadata?.name ?? ''}
        resourceVersion={machine.metadata?.resourceVersion}
      />
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
