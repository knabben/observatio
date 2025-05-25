import {ClusterInfraType, ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import {Card, Chip, Grid, GridCol, Group, Stack, Text} from "@mantine/core";
import { Pill, Table, Indicator, Space, SimpleGrid } from '@mantine/core';
import React from "react";
import {roboto, sourceCodePro400} from "@/fonts";
import Panel from "@/app/ui/dashboard/utils/panel";
import { XMarkIcon } from '@heroicons/react/24/outline';
import Specification from "@/app/ui/dashboard/components/clusters/specification";
import AITroubleshooting from "@/app/ui/dashboard/base/ai-troubleshooting";
import ObjectDetails from "@/app/ui/dashboard/base/details";
import {IconCheck, IconX} from "@tabler/icons-react";

/**
 * Displays infrastructure details of a given cluster, including cluster specifications.
 */
export default function ClusterDetails({
  cluster,
}: { cluster: ClusterType }) {
  const tabs = [
    {
      label: "Specification",
      content: (cluster: ClusterType) => <Specification cluster={cluster}/>
    },
    {
      label: "AI Troubleshooting",
      content: (cluster: ClusterType) => <AITroubleshooting conditions={cluster.status.conditions}/>
    }
  ];
  const headerRender = (cluster: ClusterType) => (
    <SimpleGrid cols={2}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            cluster.status?.infrastructureReady && cluster.status?.controlPlaneReady
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
