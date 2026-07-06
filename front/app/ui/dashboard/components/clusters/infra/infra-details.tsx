import React from "react";
import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";
import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid } from '@mantine/core';
import Specification from "@/app/ui/dashboard/components/clusters/infra/specification";
import {IconCheck, IconX} from "@tabler/icons-react";
import ObjectDetails from "@/app/ui/dashboard/base/details";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function buildContext(cluster: ClusterInfraType): ObjectContext {
  return {
    kind: 'VSphereCluster',
    name: cluster.metadata?.name ?? '',
    namespace: cluster.metadata?.namespace ?? '',
    status: cluster.status?.ready ? 'Ready' : `Not ready${cluster.status?.failureReason ? `: ${cluster.status.failureReason}` : ''}`,
    keySpecFields: {
      server: cluster.server ?? '—',
      controlPlaneEndpoint: cluster.controlPlaneEndpoint ?? '—',
    },
  };
}

/**
 * Displays infrastructure details of a given cluster, including cluster specifications,
 * vSphere cluster conditions, and associated modules.
 * It renders details in a structured layout using grid, cards, panels, and tables.
 */
export default function ClusterInfraDetails({
  cluster,
}: { cluster: ClusterInfraType }) {
  useCurrentObjectContext(buildContext(cluster));

  const tabs = [
    {
      label: "Specification",
      content: (cluster: ClusterInfraType) => <Specification cluster={cluster} />
    },
    {
      label: "YAML",
      content: (cluster: ClusterInfraType) => <ObjectTree
        gvr={RESOURCE_GVR.vsphereCluster}
        namespace={cluster.metadata?.namespace ?? ''}
        name={cluster.metadata?.name ?? ''}
        resourceVersion={cluster.metadata?.resourceVersion}
      />
    },
  ];
  const headerRender = (cluster: ClusterInfraType) => (
    <SimpleGrid cols={{base: 1, sm: 2}}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            cluster.status?.ready
              ? <IconCheck size={40} color="teal"/>
              : <IconX color="red" size={40}/>
          }
          <Text className="font-bold" fw={700}>{cluster.metadata?.name}</Text>
          <AskAIButton context={buildContext(cluster)}/>
        </Group>
      </div>
      <div>
        <Group justify="flex-end">
          <Stack gap="sm" justify="center">
            <Text size="sm">Namespace</Text>
            <Text size="xl">
              {cluster.metadata?.namespace}
            </Text>
          </Stack>
          <Stack gap="sm" justify="center">
            <Text size="sm">Age</Text>
            <Text size="xl">
              {cluster.age ?? '—'}
            </Text>
          </Stack>
        </Group>
      </div>
    </SimpleGrid>
  );
  return (
    <ObjectDetails
      object={cluster}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}
