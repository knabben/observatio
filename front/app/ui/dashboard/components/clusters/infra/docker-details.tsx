import React from "react";
import {ClusterInfraDockerType} from "@/app/ui/dashboard/components/clusters/types";
import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import AITroubleshooting from "@/app/ui/dashboard/base/ai-troubleshooting";
import {IconCheck, IconX} from "@tabler/icons-react";
import ObjectDetails from "@/app/ui/dashboard/base/details";

/**
 * Displays infrastructure details of a Docker (CAPD) cluster: readiness, namespace, age,
 * and load balancer IP, plus AI troubleshooting grounded in this resource.
 */
export default function ClusterInfraDockerDetails({
  cluster,
}: { cluster: ClusterInfraDockerType }) {
  const tabs = [
    {
      label: "AI Troubleshooting",
      content: (cluster: ClusterInfraDockerType) => <AITroubleshooting
        objectType="dockercluster"
        objectName={cluster.metadata?.name ?? ''}
        objectNamespace={cluster.metadata?.namespace ?? ''}
        conditions={[]}
      />
    }];
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
