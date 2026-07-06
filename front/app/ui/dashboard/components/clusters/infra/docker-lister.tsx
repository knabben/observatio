'use client';

import React from 'react';

import ClusterInfraDockerTable from '@/app/ui/dashboard/components/clusters/infra/docker-table'
import ClusterInfraDockerDetails from "@/app/ui/dashboard/components/clusters/infra/docker-details";

import {ClusterInfraDockerType} from "@/app/ui/dashboard/components/clusters/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

/**
 * Renders the Docker (CAPD) equivalent of `ClusterInfraLister`: live cluster infrastructure
 * details, loading/empty/error states, and search, for Docker-backed clusters.
 */
export default function ClusterInfraDockerLister() {
  return <BaseLister
    objectType="cluster-infra-docker"
    items={[]}
    renderDetails={(cluster: ClusterInfraDockerType) => <ClusterInfraDockerDetails cluster={cluster} />}
    renderTable={(clusters: ClusterInfraDockerType[], handleSelect) => (
      <ClusterInfraDockerTable select={handleSelect} clusters={clusters}/>
    )}
    title="DockerCluster / infrastructure.cluster.x-k8s.io/v1beta1"
  />
}
