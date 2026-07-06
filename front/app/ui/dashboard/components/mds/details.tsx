import React from "react";

import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";
import Specification from "@/app/ui/dashboard/components/mds/specification";

import ObjectDetails from "@/app/ui/dashboard/base/details";
import {IconCheck, IconMinus, IconX} from "@tabler/icons-react";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";

function buildContext(md: MachineDeploymentType): ObjectContext {
  const unavailable = md.status?.unavailableReplicas;
  const status = unavailable == null ? 'Unknown' : unavailable === 0 ? 'Ready' : `${unavailable} replica(s) unavailable`;
  return {
    kind: 'MachineDeployment',
    name: md.metadata?.name ?? '',
    namespace: md.metadata?.namespace ?? '',
    status,
    keySpecFields: {
      cluster: md.cluster ?? '—',
      templateVersion: md.templateversion ?? '—',
      replicas: String(md.replicas ?? '—'),
    },
  };
}

export default function MachineDeploymentDetails({
  md,
}: { md: MachineDeploymentType}) {
  useCurrentObjectContext(buildContext(md));

  const tabs = [
    {
      label: "Specification",
      content: (md: MachineDeploymentType) => <Specification md={md} />
    },
  ];
  const headerRender = (md: MachineDeploymentType) => (
    <SimpleGrid cols={{base: 1, sm: 2}}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            md.status?.unavailableReplicas == null
              ? <IconMinus color="gray" size={40}/>
              : md.status.unavailableReplicas === 0
                ? <IconCheck size={40} color="teal"/>
                : <IconX color="red" size={40}/>
          }
          <Text className="font-bold" fw={700}>{md.metadata?.name}</Text>
          <AskAIButton context={buildContext(md)}/>
        </Group>
      </div>
      <div>
        <Group justify="flex-end">
          <Stack gap="sm" justify="center">
            <Text size="sm">Namespace</Text>
            <Text size="xl">
              {md.metadata?.namespace}
            </Text>
          </Stack>
          <Stack gap="sm" justify="center">
            <Text size="sm">Age</Text>
            <Text size="xl">
              {md.age ?? '—'}
            </Text>
          </Stack>
        </Group>
      </div>
    </SimpleGrid>
  )
  return (
    <ObjectDetails
      object={md}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
