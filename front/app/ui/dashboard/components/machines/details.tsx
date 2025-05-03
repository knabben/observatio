import {Card, Grid, GridCol} from "@mantine/core";
import { Pill, Table, Indicator, Space, SimpleGrid } from '@mantine/core';
import React from "react";
import {roboto} from "@/fonts";
import Panel from "@/app/ui/dashboard/utils/panel";
import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";
import {MachineType} from "@/app/ui/dashboard/components/machines/types";

export default function MachineDetails({
  machine,
}: { machine: MachineType}) {
  return (
    <GridCol className={roboto.className} span={12}>
      <Card withBorder shadow="sm" padding="lg" radius="md">
        <SimpleGrid className="text-center" cols={3}>
          <div>
            <span className="font-bold">Name: </span>
            {
              machine.infrastructureReady && machine.bootstrapReady
              ? <Indicator offset={-3} inline withBorder position="top-end" color="green" size={10}> {machine.name} </Indicator>
              : <Indicator  offset={-3} inline withBorder position="top-end" color="red" size={10}> {machine.name} </Indicator>
            }
          </div>
          <div><span className="font-bold">Phase: </span> {machine.phase}</div>
          <div><span className="font-bold">Age:</span> {machine.created}</div>
        </SimpleGrid>
      </Card>
      <Space h="md" />
      <Grid>
        <GridCol span={12}>
          <Panel title="Specification" content={
            <Table
              variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Namespace</Table.Th>
                  <Table.Td>{machine.namespace}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Cluster</Table.Th>
                  <Table.Td>{machine.cluster}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Owner</Table.Th>
                  <Table.Td>{machine.owner}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Bootstrap</Table.Th>
                  <Table.Td>{machine.bootstrap}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Node</Table.Th>
                  <Table.Td>{machine.nodeName}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>ProviderID</Table.Th>
                  <Table.Td>{machine.providerID}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Version</Table.Th>
                  <Table.Td>{machine.version}</Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
          } />
        </GridCol>
      </Grid>
    </GridCol>
  )
}
