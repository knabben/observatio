import React from "react";
import ClusterLister from '@/app/ui/dashboard/components/clusters/lister'

export default async function Clusters() {
  return (
    <div>
      <main>
        <ClusterLister />
      </main>
    </div>
  )
}