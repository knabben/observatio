import {Card, Grid, GridCol} from "@mantine/core";
import { Pill, Table, Indicator, Space, SimpleGrid } from '@mantine/core';
import React from "react";
import {roboto, sourceCodePro400} from "@/fonts";
import Panel from "@/app/ui/dashboard/utils/panel";
import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";

export default function MachineDeploymentDetails({
  md,
}: { md: MachineDeploymentType}) {
  return (
    <GridCol className={roboto.className} span={12}>
      <Card withBorder shadow="sm" padding="lg" radius="md">
        <SimpleGrid className="text-center" cols={3}>
          <div>
            <span className="font-bold">Name: </span>
            {
              md.unavailableReplicas == 0
              ? <Indicator offset={-3} inline withBorder position="top-end" color="green" size={10}> {md.name} </Indicator>
              : <Indicator  offset={-3} inline withBorder position="top-end" color="red" size={10}> {md.name} </Indicator>
            }
          </div>
          <div><span className="font-bold">Phase: </span> {md.phase}</div>
          <div><span className="font-bold">Age:</span> {md.created}</div>
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
                  <Table.Td>{md.namespace}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Cluster</Table.Th>
                  <Table.Td>{md.cluster}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{md.replicas}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Ready Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{md.readyReplicas}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Updated Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{md.updatedReplicas}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Unavailable Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{md.unavailableReplicas}</Pill></Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
          } />
        </GridCol>
      </Grid>
    </GridCol>
  )
}
