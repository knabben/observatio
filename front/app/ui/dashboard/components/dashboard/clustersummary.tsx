'use client';

import React from 'react';
import {Card, Grid, SimpleGrid, Text, Title} from '@mantine/core';
import {getClusterSummary} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {RadialBarChart} from "@mantine/charts";
import {sourceSans400} from "@/fonts";
import { CenteredLoader } from '@/app/ui/dashboard/utils/loader';
import {useFetchState} from "@/app/ui/dashboard/shared/use-fetch-state";

type ClusterSummary = {
  name: string;
  value: number;
  color: string;
};

type ClusterSummaryResponse = {
  clusterProvisioned?: number;
  clusterFailed?: number;
  machineProvisioned?: number;
  machineFailed?: number;
  machineDeploymentProvisioned?: number;
  machineDeploymentFailed?: number;
};

const CLUSTER_STATUS_COLORS = {
  RUNNING: 'var(--mantine-color-teal-6)',
  FAILED: 'var(--mantine-color-red-6)',
} as const;

const transformSummaryData = (response: ClusterSummaryResponse): ClusterSummary[] => [
  {name: 'Cluster running', value: response.clusterProvisioned ?? 0, color: CLUSTER_STATUS_COLORS.RUNNING},
  {name: 'Cluster failed', value: response.clusterFailed ?? 0, color: CLUSTER_STATUS_COLORS.FAILED},
  {name: 'Machine D. running', value: response.machineDeploymentProvisioned ?? 0, color: CLUSTER_STATUS_COLORS.RUNNING},
  {name: 'Machine D. failed', value: response.machineDeploymentFailed ?? 0, color: CLUSTER_STATUS_COLORS.FAILED},
  {name: 'Machine provisioned', value: response.machineProvisioned ?? 0, color: CLUSTER_STATUS_COLORS.RUNNING},
  {name: 'Machine failed', value: response.machineFailed ?? 0, color: CLUSTER_STATUS_COLORS.FAILED},
];

export const useClusterSummary = () => {
  const {data: response, isLoading, error} = useFetchState<ClusterSummaryResponse | undefined>(
    getClusterSummary,
    undefined,
    'Failed to load cluster summary',
  );
  const summary = response ? transformSummaryData(response) : [];
  return {summary, isLoading, error};
}

/**
 * ClusterSummary is a React functional component that displays a summary of cluster health
 * statistics in a card layout. It fetches cluster summary data, including the number of
 * running and failed clusters, and presents the data both in textual and graphical formats.
 */
export default function ClusterSummary() {
  const {summary, isLoading, error} = useClusterSummary();

  return (
    <Card shadow="md" className={sourceSans400.className} padding="lg" radius="md" withBorder>
      <Header title="Clusters Health"/>
      {isLoading && <CenteredLoader />}
      {error && <Text c="red">{error}</Text>}
      {!error && !isLoading && (
      <Grid align="center" ta="center">
          <Grid.Col span={{base: 12, sm: 6}}>
            <SimpleGrid cols={2} verticalSpacing="sm">
              {summary.map((item: ClusterSummary, index: number) => (
                <div key={index}>
                  <div className="text-left">{item.name}<Title order={3} c={item.color}>{item.value}</Title></div>
                </div>
              ))}
            </SimpleGrid>
          </Grid.Col>
          <Grid.Col span={{base: 12, sm: 6}}>
            <RadialBarChart data={summary} dataKey="value" h={250}/>
          </Grid.Col>
        </Grid>
      )}
    </Card>
  );
}