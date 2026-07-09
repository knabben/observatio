import React from "react";

import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import {MachineSetType} from "@/app/ui/dashboard/components/machinesets/types";
import Specification from "@/app/ui/dashboard/components/machinesets/specification";

import ObjectDetails from "@/app/ui/dashboard/base/details";
import {IconCheck, IconMinus, IconX} from "@tabler/icons-react";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function msReady(desired: number | undefined, available: number | undefined): boolean | undefined {
  if (desired == null || available == null) return undefined;
  return available === desired;
}

function buildContext(ms: MachineSetType): ObjectContext {
  const ready = msReady(ms.replicas, ms.status?.availableReplicas);
  const status = ready == null ? 'Unknown' : ready ? 'Ready' : `${ms.status?.availableReplicas ?? 0}/${ms.replicas ?? 0} available`;
  return {
    kind: 'MachineSet',
    name: ms.metadata?.name ?? '',
    namespace: ms.metadata?.namespace ?? '',
    status,
    keySpecFields: {
      cluster: ms.cluster ?? '—',
      machineDeployment: ms.machineDeployment ?? '—',
      replicas: String(ms.replicas ?? '—'),
    },
  };
}

export default function MachineSetDetails({
  ms,
}: { ms: MachineSetType}) {
  useCurrentObjectContext(buildContext(ms));

  const tabs = [
    {
      label: "Specification",
      content: (ms: MachineSetType) => <Specification ms={ms} />
    },
    {
      label: "YAML",
      content: (ms: MachineSetType) => <ObjectTree
        gvr={RESOURCE_GVR.machineSet}
        namespace={ms.metadata?.namespace ?? ''}
        name={ms.metadata?.name ?? ''}
        resourceVersion={ms.metadata?.resourceVersion}
      />
    },
  ];
  const headerRender = (ms: MachineSetType) => {
    const ready = msReady(ms.replicas, ms.status?.availableReplicas);
    return (
      <SimpleGrid cols={{base: 1, sm: 2}}>
        <div className="flex items-center h-full">
          <Group justify="flex-start">
            {
              ready == null
                ? <IconMinus color="gray" size={40}/>
                : ready
                  ? <IconCheck size={40} color="teal"/>
                  : <IconX color="red" size={40}/>
            }
            <Text className="font-bold" fw={700}>{ms.metadata?.name}</Text>
            <AskAIButton context={buildContext(ms)}/>
          </Group>
        </div>
        <div>
          <Group justify="flex-end">
            <Stack gap="sm" justify="center">
              <Text size="sm">Namespace</Text>
              <Text size="xl">
                {ms.metadata?.namespace}
              </Text>
            </Stack>
            <Stack gap="sm" justify="center">
              <Text size="sm">Age</Text>
              <Text size="xl">
                {ms.age ?? '—'}
              </Text>
            </Stack>
          </Group>
        </div>
      </SimpleGrid>
    );
  }
  return (
    <ObjectDetails
      object={ms}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
