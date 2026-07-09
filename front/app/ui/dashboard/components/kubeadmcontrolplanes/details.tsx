import React from "react";

import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import {KubeadmControlPlaneType} from "@/app/ui/dashboard/components/kubeadmcontrolplanes/types";
import Specification from "@/app/ui/dashboard/components/kubeadmcontrolplanes/specification";

import ObjectDetails from "@/app/ui/dashboard/base/details";
import {IconCheck, IconMinus, IconX} from "@tabler/icons-react";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function buildContext(kcp: KubeadmControlPlaneType): ObjectContext {
  const ready = kcp.status?.ready;
  const status = ready == null ? 'Unknown' : ready ? 'Ready' : `${kcp.status?.readyReplicas ?? 0}/${kcp.status?.replicas ?? 0} ready`;
  return {
    kind: 'KubeadmControlPlane',
    name: kcp.metadata?.name ?? '',
    namespace: kcp.metadata?.namespace ?? '',
    status,
    keySpecFields: {
      cluster: kcp.cluster ?? '—',
      version: kcp.version ?? '—',
      replicas: String(kcp.replicas ?? '—'),
    },
  };
}

export default function KubeadmControlPlaneDetails({
  kcp,
}: { kcp: KubeadmControlPlaneType}) {
  useCurrentObjectContext(buildContext(kcp));

  const tabs = [
    {
      label: "Specification",
      content: (kcp: KubeadmControlPlaneType) => <Specification kcp={kcp} />
    },
    {
      label: "YAML",
      content: (kcp: KubeadmControlPlaneType) => <ObjectTree
        gvr={RESOURCE_GVR.kubeadmControlPlane}
        namespace={kcp.metadata?.namespace ?? ''}
        name={kcp.metadata?.name ?? ''}
        resourceVersion={kcp.metadata?.resourceVersion}
      />
    },
  ];
  const headerRender = (kcp: KubeadmControlPlaneType) => {
    const ready = kcp.status?.ready;
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
            <Text className="font-bold" fw={700}>{kcp.metadata?.name}</Text>
            <AskAIButton context={buildContext(kcp)}/>
          </Group>
        </div>
        <div>
          <Group justify="flex-end">
            <Stack gap="sm" justify="center">
              <Text size="sm">Namespace</Text>
              <Text size="xl">
                {kcp.metadata?.namespace}
              </Text>
            </Stack>
            <Stack gap="sm" justify="center">
              <Text size="sm">Age</Text>
              <Text size="xl">
                {kcp.age ?? '—'}
              </Text>
            </Stack>
          </Group>
        </div>
      </SimpleGrid>
    );
  }
  return (
    <ObjectDetails
      object={kcp}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
