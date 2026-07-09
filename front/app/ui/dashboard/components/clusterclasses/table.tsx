'use client';

import React from "react";
import {Badge} from '@mantine/core';
import {GridCol} from '@mantine/core';

import {ClusterClassType} from '@/app/ui/dashboard/components/clusterclasses/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {toStatusState} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

/** ClusterClass is ready when every reported condition is true; unknown when none are reported. */
function ccReady(conditions: ClusterClassType['conditions']): boolean | undefined {
  if (conditions == null || conditions.length === 0) return undefined;
  return conditions.every((c) => c.status?.toLowerCase() === 'true');
}

const columns: ColumnDef<ClusterClassType>[] = [
  {header: 'Name', render: (cc) => cc.name ?? '—'},
  {header: 'Namespace', render: (cc) => <Badge variant="light" color="gray">{cc.namespace ?? '—'}</Badge>},
  {header: 'Generation', align: 'center', render: (cc) => cc.generation ?? '—'},
  {
    header: 'Status',
    align: 'center',
    render: (cc) => <StatusIndicator state={toStatusState(ccReady(cc.conditions))} dotOnly/>,
  },
];

export default function ClusterClassTable({
  ccs, select
}: {
  ccs: ClusterClassType[]
  select: (cc: ClusterClassType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={ccs}
        columns={columns}
        getRowKey={(cc, i) => cc.metadata?.name ?? cc.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No cluster classes found"
      />
    </GridCol>
  );
}
