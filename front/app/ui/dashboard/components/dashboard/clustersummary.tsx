'use client';

import React, {useState, useEffect} from 'react';
import {Card, Grid, SimpleGrid, Text} from '@mantine/core';
import {getClusterSummary} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {RadialBarChart} from "@mantine/charts";
import {sourceCodePro400} from "@/fonts";
import { CenteredLoader } from '@/app/ui/dashboard/utils/loader';

type ClusterSummary = {
  name: string;
  value: number;
  color: string;
};

type ClusterSummaryResponse = {
  clusterProvisioned: number;
  clusterFailed: number;
  machineProvisioned: number;
  machineFailed: number;
  machineDeploymentProvisioned: number;
  machineDeploymentFailed: number;
};

const CLUSTER_STATUS_COLORS = {
  RUNNING: '#39b69d',
  FAILED: '#f53f5e',
} as const;

const transformSummaryData = (response: ClusterSummaryResponse): ClusterSummary[] => [
  {name: 'Cluster running', value: response.clusterProvisioned, color: CLUSTER_STATUS_COLORS.RUNNING},
  {name: 'Cluster failed', value: response.clusterFailed, color: CLUSTER_STATUS_COLORS.FAILED},
  {name: 'Machine D. running', value: response.machineDeploymentProvisioned, color: CLUSTER_STATUS_COLORS.RUNNING},
  {name: 'Machine D. failed', value: response.machineDeploymentFailed, color: CLUSTER_STATUS_COLORS.FAILED},
  {name: 'Machine provisioned', value: response.machineProvisioned, color: CLUSTER_STATUS_COLORS.RUNNING},
  {name: 'Machine failed', value: response.machineFailed, color: CLUSTER_STATUS_COLORS.FAILED},
];

export const useClusterSummary = () => {
  const [summary, setSummary] = useState<ClusterSummary[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleFetchError = (error: Error) => {
    console.error('Failed to fetch cluster summary:', error);
    setError('Failed to load cluster summary');
    setIsLoading(false);
  };

  useEffect(() => {
    const fetchSummary = async () => {
      try {
        const response = await getClusterSummary();
        setSummary(transformSummaryData(response));
      } catch (error) {
        handleFetchError(error as Error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchSummary();
  }, []);

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
    <Card shadow="md" className={sourceCodePro400.className} padding="lg" radius="md" withBorder>
      <Header title="Clusters Health"/>
      {isLoading && <CenteredLoader />}
      {error && <Text c="red">{error}</Text>}
      {!error && !isLoading && (
      <Grid align="center" ta="center">
          <Grid.Col span={6}>
            <SimpleGrid cols={2} verticalSpacing="sm">
              {summary.map((item: ClusterSummary, index: number) => (
                <div key={index}>
                  <div className="text-left">{item.name}</div>
                  <div >{item.value}</div>
                </div>
              ))}
            </SimpleGrid>
          </Grid.Col>
          <Grid.Col span={6}>
            <RadialBarChart data={summary} dataKey="value" h={250}/>
          </Grid.Col>
        </Grid>
      )}
    </Card>
  );
}