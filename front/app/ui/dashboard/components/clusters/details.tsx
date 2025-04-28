import {ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import {Card, Grid, GridCol} from "@mantine/core";
import { Pill, Table, Title, Indicator, Divider, Group, List } from '@mantine/core';
import React from "react";
import {sourceCodePro400} from "@/fonts";


export default function ClusterDetails({
  cluster,
}: { cluster: ClusterType }) {
  return (
    <GridCol className={sourceCodePro400.className} span={12}>
        <Group className="text-xl text-right">
          Cluster Name:
        {cluster.controlPlaneReady && cluster.infrastructureReady
          ? <Indicator offset={0} position="top-end" inline processing color="green" size={10}>{cluster.name}</Indicator>
          : <Indicator  offset={2} position="top-end" inline processing color="red" size={10}>{cluster.name}</Indicator>
        }
        <Divider size="sm" orientation="vertical" />
        Phase: {cluster.phase}
        <Divider size="sm" orientation="vertical" />
        Created: {cluster.created}
        </Group>
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Grid>
          <GridCol span={6}>
            <Title order={4}>Specification</Title>
            <Table className="text-lg" horizontalSpacing="sm" verticalSpacing="sm" variant="vertical">
              <Table.Tbody>
                <Table.Tr>
                  <Table.Th w={160}>Paused</Table.Th>
                  <Table.Td><Pill size="md">{cluster.paused.toString()}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Infrastructure Ready</Table.Th>
                  <Table.Td><Pill size="md">{cluster.infrastructureReady.toString()}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Control Plane Ready</Table.Th>
                  <Table.Td><Pill size="md">{cluster.controlPlaneReady.toString()}</Pill></Table.Td>
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
            <Title>Cluster conditions</Title>
            <List  spacing="md"
                   size="lg"
                   center>
            {
              cluster.conditions?.map((condition) => (
                <List.Item>{condition.lastTransitionTime}: {condition.type} - {condition.status}</List.Item>
              ))
            }
            </List>
        </GridCol>
        <GridCol span={6}>
          { cluster.clusterClass.isClusterClass ?
            <>
              <Title order={4}>Cluster class</Title>
              <Table horizontalSpacing="sm" verticalSpacing="sm" variant="vertical">
                <Table.Tbody>
                  <Table.Tr>
                    <Table.Th w={160}>Kubernetes Version</Table.Th>
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
            </>
            : <div />
          }
        </GridCol>
          </Grid>
      </Card>
    </GridCol>
  )
}
