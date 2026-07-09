'use client';

import React from 'react';

import KubeadmControlPlaneTable from '@/app/ui/dashboard/components/kubeadmcontrolplanes/table'
import KubeadmControlPlaneDetails from "@/app/ui/dashboard/components/kubeadmcontrolplanes/details";

import {KubeadmControlPlaneType} from "@/app/ui/dashboard/components/kubeadmcontrolplanes/types";
import BaseLister from "@/app/ui/dashboard/base/lister";
/**
 * Thin composition of `BaseLister` with the KubeadmControlPlane-specific table/details renderers.
 * `BaseLister` owns the live WebSocket stream, loading/empty/error states, and selection.
 */
export default function KubeadmControlPlaneLister() {
  return <BaseLister
    objectType="kubeadmcontrolplane"
    items={[]}
    renderDetails={(item: KubeadmControlPlaneType) => <KubeadmControlPlaneDetails kcp={item}/>}
    renderTable={(items : KubeadmControlPlaneType[], handleSelect) =>  (
      <KubeadmControlPlaneTable select={handleSelect} kcps={items}/>
    )}
    title="Kubeadm Control Planes / controlplane.cluster.x-k8s.io"
  />
}
