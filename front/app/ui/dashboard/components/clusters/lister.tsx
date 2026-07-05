'use client';

import React from 'react';

import ClusterTable from '@/app/ui/dashboard/components/clusters/table'
import ClusterDetails from "@/app/ui/dashboard/components/clusters/details";

import {ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

/**
 * Thin composition of `BaseLister` with the cluster-specific table/details renderers.
 * `BaseLister` owns the live WebSocket stream, loading/empty/error states, and selection.
 */
export default function ClusterLister() {
  return <BaseLister
    objectType="cluster"
    items={[]}
    renderDetails={(cluster: ClusterType) => <ClusterDetails cluster={cluster} />}
    renderTable={(clusters: ClusterType[], handleSelect) =>  (
      <ClusterTable select={handleSelect} clusters={clusters}/>
    )}
    title="Clusters / cluster.x-k8s.io"
  />
}