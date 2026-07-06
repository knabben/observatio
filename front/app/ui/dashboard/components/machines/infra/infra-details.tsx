import React from "react";
import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {
  Text,
  Group,
  SimpleGrid,
  Stack,
} from "@mantine/core";
import {IconCheck, IconX} from "@tabler/icons-react";
import Specification from "@/app/ui/dashboard/components/machines/infra/specification";
import ObjectDetails from "@/app/ui/dashboard/base/details";

export default function MachineInfraDetails({
  machine
}: {machine: MachineInfraType}) {
  const tabs = [
    {
      label: "Specification",
      content: (machine: MachineInfraType) => <Specification machine={machine} />
    }];

  const headerRender = (machine: MachineInfraType) => (
    <SimpleGrid cols={{base: 1, sm: 2}}>
      <div className="flex items-center h-full">
        <Group justify="flex-start">
          {
            machine.status?.ready
              ? <IconCheck size={40} color="teal"/>
              : <IconX color="red" size={40}/>
          }
          <Text className="font-bold" fw={700}>{machine.metadata?.name}</Text>
        </Group>
      </div>
      <div>
        <Group justify="flex-end">
          <Stack gap="sm" justify="center">
            <Text size="sm">Namespace</Text>
            <Text size="xl">
              {machine.metadata?.namespace}
            </Text>
          </Stack>
          <Stack gap="sm" justify="center">
            <Text size="sm">Age</Text>
            <Text size="xl">
              {machine.age ?? '—'}
            </Text>
          </Stack>
        </Group>
      </div>
    </SimpleGrid>
  );

  return (
    <ObjectDetails
      object={machine}
      headerRenderer={headerRender}
      tabs={tabs}
    />
  )
}