'use client';

import React from "react";

import MachineLister from '@/app/ui/dashboard/components/machines/lister'
import MachineInfraLister from '@/app/ui/dashboard/components/machines/infra/infra-lister'
import MachineInfraDockerLister from '@/app/ui/dashboard/components/machines/infra/docker-lister'
import {Text, Space, Tabs, TabsList, TabsTab, TabsPanel} from '@mantine/core'

import {getInfraCapabilities, emptyInfrastructureCapability} from '@/app/lib/data';
import {useFetchState} from '@/app/ui/dashboard/shared/use-fetch-state';
import {EmptyState} from '@/app/ui/dashboard/shared/empty-state';
import {CenteredLoader} from '@/app/ui/dashboard/utils/loader';

/**
 * Renders the Machines screen's tabs, adapting to whichever infrastructure provider(s) are
 * actually installed, mirroring `ClusterTabs` for the Clusters screen.
 */
export default function MachineTabs() {
  const {data: capability, isLoading} = useFetchState(
    getInfraCapabilities,
    emptyInfrastructureCapability,
    'Failed to load infrastructure capabilities',
  );

  if (isLoading) {
    return <CenteredLoader/>;
  }

  const noProviderDetected = !capability.docker.installed && !capability.vsphere.installed;

  return (
    <Tabs color="var(--mantine-color-brand-4)" defaultValue="machine">
      <TabsList>
        <TabsTab value="machine">
          <Text size="md" fw={700}>Machines</Text>
        </TabsTab>
        {capability.docker.installed &&
          <TabsTab value="machine-docker">
            <Text size="md" fw={700}>Docker Machines</Text>
          </TabsTab>
        }
        {capability.vsphere.installed &&
          <TabsTab value="machine-vsphere">
            <Text size="md" fw={700}>vSphere Machines</Text>
          </TabsTab>
        }
      </TabsList>

      <TabsPanel value="machine">
        <Space h="lg"/>
        <MachineLister/>
        {noProviderDetected &&
          <>
            <Space h="lg"/>
            <EmptyState label="No supported infrastructure provider detected"/>
          </>
        }
      </TabsPanel>

      {capability.docker.installed &&
        <TabsPanel value="machine-docker">
          <Space h="lg"/>
          <MachineInfraDockerLister/>
        </TabsPanel>
      }

      {capability.vsphere.installed &&
        <TabsPanel value="machine-vsphere">
          <Space h="lg"/>
          <MachineInfraLister/>
        </TabsPanel>
      }
    </Tabs>
  )
}
