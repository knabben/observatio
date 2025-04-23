import Link from 'next/link';

import ClusterLister from '@/app/ui/dashboard/components/ClusterLister'
import Search from "@/app/ui/dashboard/search";
import React from "react";

export default async function Clusters() {
  return (
    <div>
      <main>
        <ClusterLister />
      </main>
    </div>
  )
}