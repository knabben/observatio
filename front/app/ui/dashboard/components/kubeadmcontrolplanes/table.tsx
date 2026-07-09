'use client';

import React from "react";
import {Badge} from '@mantine/core';
import {GridCol} from '@mantine/core';

import {KubeadmControlPlaneType} from '@/app/ui/dashboard/components/kubeadmcontrolplanes/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {toStatusState} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

const columns: ColumnDef<KubeadmControlPlaneType>[] = [
  {header: 'Name', render: (kcp) => kcp.metadata?.name ?? '—'},
  {header: 'Namespace', render: (kcp) => <Badge variant="light" color="gray">{kcp.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Cluster', render: (kcp) => kcp.cluster ?? '—'},
  {header: 'Version', render: (kcp) => kcp.version ?? '—'},
  {header: 'Ready / Replicas', align: 'center', render: (kcp) => `${kcp.status?.readyReplicas ?? '—'} / ${kcp.status?.replicas ?? '—'}`},
  {header: 'Age', align: 'center', render: (kcp) => kcp.age ?? '—'},
  {
    header: 'Status',
    align: 'center',
    render: (kcp) => <StatusIndicator state={toStatusState(kcp.status?.ready)} dotOnly/>,
  },
];

export default function KubeadmControlPlaneTable({
  kcps, select
}: {
  kcps: KubeadmControlPlaneType[]
  select: (kcp: KubeadmControlPlaneType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={kcps}
        columns={columns}
        getRowKey={(kcp, i) => kcp.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No kubeadm control planes found"
      />
    </GridCol>
  );
}
