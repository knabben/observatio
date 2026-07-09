import React from "react";

import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import {ClusterClassType} from "@/app/ui/dashboard/components/clusterclasses/types";
import Specification from "@/app/ui/dashboard/components/clusterclasses/specification";

import ObjectDetails from "@/app/ui/dashboard/base/details";
import {IconCheck, IconMinus, IconX} from "@tabler/icons-react";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function ccReady(conditions: ClusterClassType['conditions']): boolean | undefined {
  if (conditions == null || conditions.length === 0) return undefined;
  return conditions.every((c) => c.status?.toLowerCase() === 'true');
}

function buildContext(cc: ClusterClassType): ObjectContext {
  const ready = ccReady(cc.conditions);
  const status = ready == null ? 'Unknown' : ready ? 'Ready' : 'Degraded';
  return {
    kind: 'ClusterClass',
    name: cc.name ?? '',
    namespace: cc.namespace ?? '',
    status,
    keySpecFields: {
      generation: String(cc.generation ?? '—'),
    },
  };
}

export default function ClusterClassDetails({
  cc,
}: { cc: ClusterClassType}) {
  useCurrentObjectContext(buildContext(cc));

  const tabs = [
    {
      label: "Specification",
      content: (cc: ClusterClassType) => <Specification cc={cc} />
    },
    {
      label: "YAML",
      content: (cc: ClusterClassType) => <ObjectTree
        gvr={RESOURCE_GVR.clusterClass}
        namespace={cc.namespace ?? ''}
        name={cc.name ?? ''}
      />
    },
  ];
  const headerRender = (cc: ClusterClassType) => {
    const ready = ccReady(cc.conditions);
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
            <Text className="font-bold" fw={700}>{cc.name}</Text>
            <AskAIButton context={buildContext(cc)}/>
          </Group>
        </div>
        <div>
          <Group justify="flex-end">
            <Stack gap="sm" justify="center">
              <Text size="sm">Namespace</Text>
              <Text size="xl">
                {cc.namespace}
              </Text>
            </Stack>
          </Group>
        </div>
      </SimpleGrid>
    );
  }
  return (
    <ObjectDetails
      object={cc}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
