import React from "react";

import ClusterLister from '@/app/ui/dashboard/components/clusters/lister'
import { Text, Space, Tabs, TabsList, TabsTab, TabsPanel } from '@mantine/core'

export default async function Clusters() {
  return (
      <Tabs color="#aaf16a" defaultValue="clusters">
        <TabsList>
          <TabsTab value="clusters">
            <Text size="md" fw={700}>Clusters</Text>
          </TabsTab>
          <TabsTab value="vsphere">
            <Text size="md" fw={700}> vSphere Clusters</Text>
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