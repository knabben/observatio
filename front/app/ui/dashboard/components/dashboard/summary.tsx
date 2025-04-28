'use client';

import React, {useState, useEffect} from 'react';
import {Card, Grid, Space, Text} from '@mantine/core';
import {getClusterSummary} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {sourceCodePro400} from "@/fonts";

type Summary = {
  failed?: bigint
  provisioned?: bigint,
}

// Summary : Resume clusters statuses.
export default function ClusterSummary() {
  const [clusterSummary, setClusterSummary] = useState<Summary>({});

  useEffect( () => {
    const fetch = async  () => {
      setClusterSummary(await getClusterSummary())
    }
    fetch().catch( (e) => { console.error('error', e) })
  }, [])

  return (
    <Card shadow="md" className={sourceCodePro400.className} padding="lg" radius="md" withBorder>
      <Header title="Clusters Health"/>
      <div style={{ resize: 'vertical', overflow: 'hidden', maxHeight: '100%' }}>
        <Grid ta="center">
          <Grid.Col span={6}>
            <Text size="lg">Provisioned</Text>
            <Space h="lg" />
            <div className="text-5xl text-[#39b69d]">{clusterSummary.provisioned}</div>
          </Grid.Col>
          <Grid.Col span={6}>
            Failing
            <Space h="lg"/>
            <div className="text-5xl text-[#e01003]">{clusterSummary.failed}</div>
          </Grid.Col>
        </Grid>
      </div>
    </Card>
  );
}