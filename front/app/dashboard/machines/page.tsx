import React from "react";
import MachineLister from '@/app/ui/dashboard/components/machines/lister'
import MachineInfraLister from '@/app/ui/dashboard/components/machines/infra-lister'
import { Text, Space, Tabs, TabsList, TabsTab, TabsPanel } from '@mantine/core'

export default async function Machines() {
  return (
    <Tabs color="#aaf16a" defaultValue="machine-vsphere">
      <TabsList>
        <TabsTab value="machine-vsphere">
          <Text size="md" fw={700}>vSphere Machines</Text>
        </TabsTab>
        <TabsTab value="machine">
          <Text size="md" fw={700}>Machines</Text>
        </TabsTab>
      </TabsList>
      <TabsPanel value="machine">
        <Space h="lg" />
        <MachineLister />
      </TabsPanel>
      <TabsPanel value="machine-vsphere">
        <Space h="lg" />
        <MachineInfraLister />
      </TabsPanel>
    </Tabs>
  )
}
