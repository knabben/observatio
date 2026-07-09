import React from "react";
import KubeadmControlPlaneLister from '@/app/ui/dashboard/components/kubeadmcontrolplanes/lister'

export default async function KubeadmControlPlanes() {
  return (
    <div>
      <main>
        <KubeadmControlPlaneLister />
      </main>
    </div>
  )
}
