import React from "react";
import {ClusterInfraDockerType} from "@/app/ui/dashboard/components/clusters/types";
import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import DockerSpecification from "@/app/ui/dashboard/components/clusters/infra/docker-specification";
import {IconCheck, IconX} from "@tabler/icons-react";
import ObjectDetails from "@/app/ui/dashboard/base/details";
import {ObjectContext} from "@/app/ui/dashboard/ai-panel/ai-panel-context";
import {useCurrentObjectContext} from "@/app/ui/dashboard/ai-panel/use-current-object-context";
import {AskAIButton} from "@/app/ui/dashboard/ai-panel/ask-ai-button";
import {ObjectTree} from "@/app/ui/dashboard/shared/object-tree";
import {RESOURCE_GVR} from "@/app/lib/resource-gvr";

function buildContext(cluster: ClusterInfraDockerType): ObjectContext {
  return {
    kind: 'DockerCluster',
    name: cluster.metadata?.name ?? '',
    namespace: cluster.metadata?.namespace ?? '',
    status: cluster.ready ? 'Ready' : 'Not ready',
    keySpecFields: {
      loadBalancerIP: cluster.loadBalancerIP ?? '—',
    },
  };
}

/**
 * Displays infrastructure details of a Docker (CAPD) cluster: readiness, namespace, age,
 * and load balancer IP.
 */
export default function ClusterInfraDockerDetails({
  cluster,
}: { cluster: ClusterInfraDockerType }) {
  useCurrentObjectContext(buildContext(cluster));

  const tabs = [
    {
      label: "Specification",
      content: (cluster: ClusterInfraDockerType) => <DockerSpecification cluster={cluster} />
    },
    {
      label: "YAML",
      content: (cluster: ClusterInfraDockerType) => <ObjectTree
        gvr={RESOURCE_GVR.dockerCluster}
        namespace={cluster.metadata?.namespace ?? ''}
        name={cluster.metadata?.name ?? ''}
        resourceVersion={cluster.metadata?.resourceVersion}
      />
    },
  ];
  const headerRender = (cluster: ClusterInfraDockerType) => (
    <SimpleGrid cols={{base: 1, sm: 2}}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            cluster.ready
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
            <Text size="sm">Load Balancer IP</Text>
            <Text size="xl">
              {cluster.loadBalancerIP ?? '—'}
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
