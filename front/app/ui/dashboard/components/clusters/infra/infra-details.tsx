import React from "react";
import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";
import {Card, Chip, Grid, GridCol} from "@mantine/core";
import { Pill, Table, Indicator, Space, SimpleGrid } from '@mantine/core';
import {roboto, sourceCodePro400} from "@/fonts";
import Panel from "@/app/ui/dashboard/utils/panel";
import { XMarkIcon } from '@heroicons/react/24/outline';

/**
 * Displays infrastructure details of a given cluster, including cluster specifications,
 * vSphere cluster conditions, and associated modules.
 * It renders details in a structured layout using grid, cards, panels, and tables.
 */
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
        <Space h="md" />
        <GridCol span={6}>
          <Panel title="vSphere Cluster Conditions" content={
            <Table variant="vertical">
              <Table.Tbody className="text-sm">
                {
                  cluster.conditions?.map((condition,ic) => (
                    <Table.Tr key={condition.type}>
                      <Table.Td>
                        {
                          condition.status.toLowerCase() == "true"
                          ? <Chip key={ic} className="p-1" defaultChecked color="teal" variant="light">{condition.type}</Chip>
                          : <Chip key={ic} defaultChecked icon={<XMarkIcon />} color="red" variant="light">{condition.type}</Chip>
                        }
                      </Table.Td>
                      <Table.Td>{condition.lastTransitionTime}</Table.Td>
                    </Table.Tr>
                  ))
                }
              </Table.Tbody>
            </Table>
          }></Panel>
        </GridCol>
      </Grid>
    </GridCol>
  )
}
