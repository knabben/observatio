'use client';

import React, {useState, useEffect} from 'react';
import {Card, Grid} from '@mantine/core';
import {getClusterSummary} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {sourceCodePro400} from "@/fonts";
import {RadialBarChart, PieChart} from "@mantine/charts";


type Summary = {
  name: string,
  value: number,
  color: string,
}

/**
 * ClusterSummary is a React functional component that displays a summary of cluster health
 * statistics in a card layout. It fetches cluster summary data, including the number of
 * running and failed clusters, and presents the data both in textual and graphical formats.
 */
export default function ClusterSummary() {
  const [summary, setSummary] = useState<Summary[]>([]);
  useEffect( () => {
    const fetch = async  () => {
      const response = await getClusterSummary()
      setSummary([
        { name: 'Running', value: response.provisioned, color: '#39b69d' },
        { name: 'Failed', value: response.failed, color: '#f53f5e' },
      ])
    }
    fetch().catch( (e) => { console.error('error', e) })
  }, [])

  return (
    <Card shadow="md" className={sourceCodePro400.className} padding="lg" radius="md" withBorder>
      <Header title="Clusters Health"/>
      <Grid justify="center" align="center" ta="center">
        <Grid.Col span={6}>
          {
            summary?.map((summary: Summary, i: number) => (
              <div key={i}>{summary.value} {summary.name}</div>
            ))
           }
        </Grid.Col>
        <Grid.Col span={6}>
          <RadialBarChart data={summary} dataKey="value" h={250} />
        </Grid.Col>
      </Grid>
    </Card>
  );
}