'use client';

import React from 'react';

import ClusterClassTable from '@/app/ui/dashboard/components/clusterclasses/table'
import ClusterClassDetails from "@/app/ui/dashboard/components/clusterclasses/details";

import {ClusterClassType} from "@/app/ui/dashboard/components/clusterclasses/types";
import BaseLister from "@/app/ui/dashboard/base/lister";
/**
 * Thin composition of `BaseLister` with the ClusterClass-specific table/details renderers.
 * `BaseLister` owns the live WebSocket stream, loading/empty/error states, and selection. This is
 * additive to, and does not replace, the existing main-dashboard `ClusterClassLister` widget
 * (research.md R5).
 */
export default function ClusterClassLister() {
  return <BaseLister
    objectType="clusterclass"
    items={[]}
    renderDetails={(item: ClusterClassType) => <ClusterClassDetails cc={item}/>}
    renderTable={(items : ClusterClassType[], handleSelect) =>  (
      <ClusterClassTable select={handleSelect} ccs={items}/>
    )}
    title="Cluster Classes / cluster.x-k8s.io"
  />
}
