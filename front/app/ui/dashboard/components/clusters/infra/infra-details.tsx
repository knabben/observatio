import React from "react";
import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";
import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid } from '@mantine/core';
import Specification from "@/app/ui/dashboard/components/clusters/infra/specification";
import AITroubleshooting from "@/app/ui/dashboard/base/ai-troubleshooting";
import {IconCheck, IconX} from "@tabler/icons-react";
import ObjectDetails from "@/app/ui/dashboard/base/details";

/**
 * Displays infrastructure details of a given cluster, including cluster specifications,
 * vSphere cluster conditions, and associated modules.
 * It renders details in a structured layout using grid, cards, panels, and tables.
 */
export default function ClusterInfraDetails({
  cluster,
}: { cluster: ClusterInfraType }) {
  const tabs = [
    {
      label: "Specification",
      content: (cluster: ClusterInfraType) => <Specification cluster={cluster} />
    },
    {
      label: "AI Troubleshooting",
      content: (cluster: ClusterInfraType) => <AITroubleshooting conditions={cluster.status.conditions} />
  }];
  const headerRender = (cluster: ClusterInfraType) => (
    <SimpleGrid cols={2}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            cluster.status?.ready
              ? <IconCheck size={40} color="teal"/>
              : <IconX color="red" size={40}/>
          }
          <Text className="text-bold" fw={700}>{cluster.metadata?.name}</Text>
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
            <Text size="sm">Created</Text>
            <Text size="xl">
              {cluster.age}
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
