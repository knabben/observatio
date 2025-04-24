'use client';

import React, {useState, useEffect} from 'react';
import { Grid, Space, Text, Divider } from '@mantine/core';
import {getClusterSummary} from "@/app/lib/data";

type Summary = {
  failed: bigint
  provisioned: bigint,
}

// Summary : Resume clusters statuses.
export default function ClusterSummary() {
  const [clusterSummary, setClusterSummary] = useState<Summary>([])
  useEffect( () => {
    const fetch = async  () => {
      setClusterSummary(await getClusterSummary())
    }
    fetch().catch( (e) => { console.error('error', e) })
  }, [])

  return (
    <div  style={{ resize: 'vertical', overflow: 'hidden', maxHeight: '100%' }}>
        <Grid ta="center">
          <Grid.Col span={6}>
            Provisioned
            <Space h="lg" />
            <Text fw={900} size="xl"  c="teal.4">{clusterSummary.provisioned}</Text>
          </Grid.Col>
          <Grid.Col span={6}>
            Failing
            <Space h="lg" />
            <Text fw={900} size="xl" c="red">{clusterSummary.failed}</Text>
          </Grid.Col>
        </Grid>
    </div>
  );
}