'use client';

import React from "react";

import {Badge} from '@mantine/core';
import {GridCol} from '@mantine/core';

import {ClusterInfraDockerType} from '@/app/ui/dashboard/components/clusters/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {toStatusState} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

const columns: ColumnDef<ClusterInfraDockerType>[] = [
  {header: 'Name', render: (c) => c.metadata?.name ?? '—'},
  {header: 'Namespace', render: (c) => <Badge variant="light" color="gray">{c.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Cluster', render: (c) => c.cluster ?? '—'},
  {header: 'Load Balancer IP', render: (c) => c.loadBalancerIP ?? '—'},
  {header: 'Age', render: (c) => c.age ?? '—'},
  {header: 'Status', align: 'center', render: (c) => <StatusIndicator state={toStatusState(c.ready)} dotOnly/>},
];

export default function ClusterInfraDockerTable({
  clusters, select
}: {
  clusters: ClusterInfraDockerType[]
  select: (cluster: ClusterInfraDockerType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={clusters}
        columns={columns}
        getRowKey={(c, i) => c.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No Docker clusters found"
      />
    </GridCol>
  )
}
