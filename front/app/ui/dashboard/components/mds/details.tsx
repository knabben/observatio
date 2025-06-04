import React from "react";

import {Group, Stack, Text} from "@mantine/core";
import {SimpleGrid} from '@mantine/core';
import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";
import Specification from "@/app/ui/dashboard/components/mds/specification";
import AITroubleshooting from "@/app/ui/dashboard/base/ai-troubleshooting";

import ObjectDetails from "@/app/ui/dashboard/base/details";
import {IconCheck, IconX} from "@tabler/icons-react";

export default function MachineDeploymentDetails({
  md,
}: { md: MachineDeploymentType}) {
  const tabs = [
    {
      label: "Specification",
      content: (md: MachineDeploymentType) => <Specification md={md} />
    },
    {
      label: "AI Troubleshooting",
      content: (md: MachineDeploymentType) => <AITroubleshooting
        objectType="machinedeployment"
        objectName={md.metadata.name}
        objectNamespace={md.metadata.namespace}
        conditions={md.status.conditions} />
    }
  ];
  const headerRender = (md: MachineDeploymentType) => (
    <SimpleGrid cols={2}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            md.status.unavailableReplicas == 0
              ? <IconCheck size={40} color="teal"/>
              : <IconX color="red" size={40}/>
          }
          <Text className="text-bold" fw={700}>{md.metadata?.name}</Text>
        </Group>
      </div>
      <div>
        <Group justify="flex-end">
          <Stack gap="sm" justify="center">
            <Text size="sm">Namespace</Text>
            <Text size="xl">
              {md.metadata?.namespace}
            </Text>
          </Stack>
          <Stack gap="sm" justify="center">
            <Text size="sm">Created</Text>
            <Text size="xl">
              {md.age}
            </Text>
          </Stack>
        </Group>
      </div>
    </SimpleGrid>
  )
  return (
    <ObjectDetails
      object={md}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}