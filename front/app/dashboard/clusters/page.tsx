import React from "react";

import ClusterLister from '@/app/ui/dashboard/components/clusters/lister'
import { Space, Tabs, TabsList, TabsTab, TabsPanel } from '@mantine/core'

export default async function Clusters() {
  return (
      <Tabs defaultValue="clusters">
        <TabsList>
          <TabsTab value="clusters">
            Clusters
          </TabsTab>
          <TabsTab value="vsphere">
            vSphere Clusters
          </TabsTab>
        </TabsList>

        <TabsPanel value="clusters">
          <Space h="lg" />
          <ClusterLister />
        </TabsPanel>

        <TabsPanel value="vsphere">
          Infra structure - vsphere
        </TabsPanel>
      </Tabs>
  )
}