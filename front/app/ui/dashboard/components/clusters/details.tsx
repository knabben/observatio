import {ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import {Card, Grid, GridCol} from "@mantine/core";
import { Pill, Table, Indicator, Space, SimpleGrid } from '@mantine/core';
import React from "react";
import {roboto, sourceCodePro400} from "@/fonts";
import Header from "@/app/ui/dashboard/utils/header";

export default function ClusterDetails({
  cluster,
}: { cluster: ClusterType }) {
  return (
    <GridCol className={roboto.className} span={12}>
      <Card withBorder shadow="sm" padding="lg" radius="md">
        <SimpleGrid className="text-center" cols={3}>
          <div>
            <span className="font-bold">Cluster Name: </span>
            {
              cluster.controlPlaneReady && cluster.infrastructureReady
              ? <Indicator offset={-3} inline withBorder position="top-end" color="green" size={7}> {cluster.name} </Indicator>
              : <Indicator  offset={-3} inline withBorder position="top-end" color="red" size={7}> {cluster.name} </Indicator>
            }
          </div>
          <div><span className="font-bold">Phase: </span> {cluster.phase}</div>
          <div><span className="font-bold">Age:</span> {cluster.created}</div>
        </SimpleGrid>
      </Card>

      <Space h="lg" />
      <Grid>
        <GridCol span={6}>
          <Card className={roboto.className} shadow="sm" padding="lg" radius="md" withBorder>
            <Header title="Specification" />
            <Table
              variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Paused</Table.Th>
                  <Table.Td><Pill size="sm">{cluster.paused.toString()}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Infrastructure Ready</Table.Th>
                  <Table.Td><Pill size="sm">{cluster.infrastructureReady.toString()}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Control Plane Ready</Table.Th>
                  <Table.Td><Pill size="sm">{cluster.controlPlaneReady.toString()}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Pod Network</Table.Th>
                  <Table.Td>{cluster.podNetwork}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Service Network</Table.Th>
                  <Table.Td>{cluster.serviceNetwork}</Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
            <Space h="lg" />
            <Header title="Cluster conditions" />
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
          </Card>
        </GridCol>
        <GridCol span={6}>
          <Card shadow="sm" padding="lg" radius="md" withBorder>
            <Header title="Cluster Class" />
            <Table variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th className="text-medium" w={230}>Kubernetes Version</Table.Th>
                  <Table.Td>{cluster.clusterClass.kubernetesVersion}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Cluster Class Name</Table.Th>
                  <Table.Td>{cluster.clusterClass.className}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Control Plane Ready</Table.Th>
                  <Table.Td>{cluster.controlPlaneReady.toString()}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Pod Network</Table.Th>
                  <Table.Td>{cluster.podNetwork}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Service Network</Table.Th>
                  <Table.Td>{cluster.serviceNetwork}</Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
            <Space h="xl" />
            <Header title="Machine Deployments" />
            <Table horizontalSpacing="sm" verticalSpacing="sm">
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>Name</Table.Th>
                  <Table.Th>Class</Table.Th>
                  <Table.Th>Replicas</Table.Th>
                  <Table.Th>Strategy</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody className="text-base">
              {
                cluster.clusterClass?.machineDeployments?.map((md) => (
                  <Table.Tr className={sourceCodePro400.className} key={cluster.name}>
                    <Table.Td>{md.name}</Table.Td>
                    <Table.Td>{md.class}</Table.Td>
                    <Table.Td>{md.replicas}</Table.Td>
                    <Table.Td>{md.strategy?.type}</Table.Td>
                  </Table.Tr>
                ))
              }
              </Table.Tbody>
            </Table>
          </Card>
        </GridCol>
      </Grid>
    </GridCol>
  )
}
