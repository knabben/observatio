'use client';

import React from 'react';

import ClusterInfraTable from '@/app/ui/dashboard/components/clusters/infra/infra-table'
import ClusterInfraDetails from "@/app/ui/dashboard/components/clusters/infra/infra-details";

import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

/**
 * The `ClusterInfraLister` function is a React functional component responsible for rendering
 * a user interface to manage and display vSphere cluster infrastructure details.
 * It handles WebSocket communication, loading states, search functionality, and conditionally
 * displays a detailed view or a table of vSphere clusters based on the data provided.
 */
export default function ClusterInfraLister() {
  return <BaseLister
    objectType="cluster-infra"
    items={[]}
    renderDetails={(cluster: ClusterInfraType) => <ClusterInfraDetails cluster={cluster} />}
    renderTable={(clusters: ClusterInfraType[], handleSelect) =>  (
      <ClusterInfraTable select={handleSelect} clusters={clusters}/>
    )}
    title="VSphereCluster / infrastructure.cluster.x-k8s.io/v1beta1"
  />
}