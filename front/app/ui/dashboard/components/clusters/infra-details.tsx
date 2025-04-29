import React from "react";
import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";
import {Card, Grid, GridCol} from "@mantine/core";
import { Pill, Table, Indicator, Space, SimpleGrid } from '@mantine/core';
import {roboto, sourceCodePro400} from "@/fonts";
import Panel from "@/app/ui/dashboard/utils/panel";

export default function ClusterInfraDetails({
  cluster,
}: { cluster: ClusterInfraType }) {
  return (
    <GridCol className={roboto.className} span={12}>
      <Card withBorder shadow="sm" padding="lg" radius="md">
        <SimpleGrid className="text-center" cols={2}>
          <div>
            <span className="font-bold">Cluster Name: </span>
            {
              cluster.ready
              ? <Indicator offset={-3} inline withBorder position="top-end" color="green" size={10}> {cluster.name} </Indicator>
              : <Indicator  offset={-3} inline withBorder position="top-end" color="red" size={10}> {cluster.name} </Indicator>
            }
          </div>
          <div><span className="font-bold">Age:</span> {cluster.created}</div>
        </SimpleGrid>
      </Card>
      <Space h="md" />
      <Grid>
        <GridCol span={6}>
          <Panel title="Specification" content={
            <Table
              variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Cluster</Table.Th>
                  <Table.Td>{cluster.cluster}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th w={260}>Control Plane Endpoint</Table.Th>
                  <Table.Td>{cluster.controlPlaneEndpoint}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th w={260}>Server</Table.Th>
                  <Table.Td>{cluster.server}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th w={260}>Thumbprint</Table.Th>
                  <Table.Td>{cluster.thumbprint}</Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
          }/>
          <Space h="md" />
          <Panel title="vSphere Cluster Conditions" content={
              <Table variant="vertical">
                <Table.Tbody className="text-sm">
                  {
                    cluster.conditions?.map((condition) => (
                      <Table.Tr key={condition.type}>
                        <Table.Th>{condition.type}</Table.Th>
                        <Table.Td><Pill size="sm">{condition.status}</Pill></Table.Td>
                        <Table.Td>{condition.lastTransitionTime}</Table.Td>
                      </Table.Tr>
                    ))
                  }
                </Table.Tbody>
              </Table>
          }></Panel>
        </GridCol>
        <Space h="md" />
        <GridCol span={6}>
          <Panel title="Cluster Modules" content={
            <Table horizontalSpacing="sm" verticalSpacing="sm">
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>Control Plane</Table.Th>
                  <Table.Th>Target Object Name</Table.Th>
                  <Table.Th>Module UUID</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody className="text-base">
                {
                  cluster.modules?.map((module) => (
                    <Table.Tr className={sourceCodePro400.className} key={module.moduleUUID}>
                      <Table.Td><Pill>{module.controlPlane.toString()}</Pill></Table.Td>
                      <Table.Td>{module.targetObjectName}</Table.Td>
                      <Table.Td>{module.moduleUUID}</Table.Td>
                    </Table.Tr>
                  ))
                }
              </Table.Tbody>
            </Table>
          }/>
        </GridCol>
      </Grid>
    </GridCol>
  )
}
