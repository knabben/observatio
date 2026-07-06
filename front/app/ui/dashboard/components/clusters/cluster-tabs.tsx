'use client';

import React from "react";

import ClusterLister from '@/app/ui/dashboard/components/clusters/lister'
import ClusterInfraLister from '@/app/ui/dashboard/components/clusters/infra/infra-lister'
import ClusterInfraDockerLister from '@/app/ui/dashboard/components/clusters/infra/docker-lister'
import {Text, Space, Tabs, TabsList, TabsTab, TabsPanel} from '@mantine/core'

import {getInfraCapabilities, emptyInfrastructureCapability} from '@/app/lib/data';
import {useFetchState} from '@/app/ui/dashboard/shared/use-fetch-state';
import {EmptyState} from '@/app/ui/dashboard/shared/empty-state';
import {CenteredLoader} from '@/app/ui/dashboard/utils/loader';

/**
 * Renders the Clusters screen's tabs, adapting to whichever infrastructure provider(s) are
 * actually installed in the connected environment instead of always assuming vSphere: a
 * provider's tab only appears when detected, both appear in a mixed environment, and a clear
 * message is shown when neither is detected.
 */
export default function ClusterTabs() {
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
    <Tabs color="var(--mantine-color-brand-4)" defaultValue="clusters">
      <TabsList>
        <TabsTab value="clusters">
          <Text size="md" fw={700}>Clusters</Text>
        </TabsTab>
        {capability.docker.installed &&
          <TabsTab value="docker">
            <Text size="md" fw={700}>Docker Clusters</Text>
          </TabsTab>
        }
        {capability.vsphere.installed &&
          <TabsTab value="vsphere">
            <Text size="md" fw={700}>vSphere Clusters</Text>
          </TabsTab>
        }
      </TabsList>

      <TabsPanel value="clusters">
        <Space h="lg"/>
        <ClusterLister/>
        {noProviderDetected &&
          <>
            <Space h="lg"/>
            <EmptyState label="No supported infrastructure provider detected"/>
          </>
        }
      </TabsPanel>

      {capability.docker.installed &&
        <TabsPanel value="docker">
          <Space h="lg"/>
          <ClusterInfraDockerLister/>
        </TabsPanel>
      }

      {capability.vsphere.installed &&
        <TabsPanel value="vsphere">
          <Space h="lg"/>
          <ClusterInfraLister/>
        </TabsPanel>
      }
    </Tabs>
  )
}
