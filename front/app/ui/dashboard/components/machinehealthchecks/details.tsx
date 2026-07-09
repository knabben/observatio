import React from "react";

import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import {MachineHealthCheckType} from "@/app/ui/dashboard/components/machinehealthchecks/types";
import Specification from "@/app/ui/dashboard/components/machinehealthchecks/specification";

import ObjectDetails from "@/app/ui/dashboard/base/details";
import {IconCheck, IconMinus, IconX} from "@tabler/icons-react";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function mhcReady(expected: number | undefined, healthy: number | undefined): boolean | undefined {
  if (expected == null || healthy == null) return undefined;
  return healthy === expected;
}

function buildContext(mhc: MachineHealthCheckType): ObjectContext {
  const ready = mhcReady(mhc.status?.expectedMachines, mhc.status?.currentHealthy);
  const status = ready == null ? 'Unknown' : ready ? 'Ready' : `${mhc.status?.currentHealthy ?? 0}/${mhc.status?.expectedMachines ?? 0} healthy`;
  return {
    kind: 'MachineHealthCheck',
    name: mhc.metadata?.name ?? '',
    namespace: mhc.metadata?.namespace ?? '',
    status,
    keySpecFields: {
      cluster: mhc.cluster ?? '—',
      maxUnhealthy: mhc.maxUnhealthy ?? '—',
      expectedMachines: String(mhc.status?.expectedMachines ?? '—'),
    },
  };
}

export default function MachineHealthCheckDetails({
  mhc,
}: { mhc: MachineHealthCheckType}) {
  useCurrentObjectContext(buildContext(mhc));

  const tabs = [
    {
      label: "Specification",
      content: (mhc: MachineHealthCheckType) => <Specification mhc={mhc} />
    },
    {
      label: "YAML",
      content: (mhc: MachineHealthCheckType) => <ObjectTree
        gvr={RESOURCE_GVR.machineHealthCheck}
        namespace={mhc.metadata?.namespace ?? ''}
        name={mhc.metadata?.name ?? ''}
        resourceVersion={mhc.metadata?.resourceVersion}
      />
    },
  ];
  const headerRender = (mhc: MachineHealthCheckType) => {
    const ready = mhcReady(mhc.status?.expectedMachines, mhc.status?.currentHealthy);
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
            <Text className="font-bold" fw={700}>{mhc.metadata?.name}</Text>
            <AskAIButton context={buildContext(mhc)}/>
          </Group>
        </div>
        <div>
          <Group justify="flex-end">
            <Stack gap="sm" justify="center">
              <Text size="sm">Namespace</Text>
              <Text size="xl">
                {mhc.metadata?.namespace}
              </Text>
            </Stack>
            <Stack gap="sm" justify="center">
              <Text size="sm">Age</Text>
              <Text size="xl">
                {mhc.age ?? '—'}
              </Text>
            </Stack>
          </Group>
        </div>
      </SimpleGrid>
    );
  }
  return (
    <ObjectDetails
      object={mhc}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
