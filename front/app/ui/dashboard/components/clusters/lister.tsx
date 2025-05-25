'use client';

import React from 'react';

import ClusterTable from '@/app/ui/dashboard/components/clusters/table'
import ClusterDetails from "@/app/ui/dashboard/components/clusters/details";

import {ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

/**
 * A functional component that fetches, filters, and displays a list of clusters.
 * The component integrates WebSocket for real-time communication and enables
 * cluster search functionality. Displays a loader while data is being fetched.
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