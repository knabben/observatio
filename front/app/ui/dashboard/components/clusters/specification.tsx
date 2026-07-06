import {ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import {Grid, GridCol, Pill, Space, Table} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import ConditionsTable from "@/app/ui/dashboard/shared/conditions-table";
import {sourceSans400} from "@/fonts";
import React from "react";

export default function Specification({
  cluster,
}: {cluster: ClusterType}) {
  return (
    <Grid>
      <GridCol span={{base: 12, md: 6}}>
        <Panel title="Specification" content={
          <Table
            variant="vertical">
            <Table.Tbody className="text-sm">
              <Table.Tr>
                <Table.Th w={260}>Paused</Table.Th>
                <Table.Td><Pill size="sm">{cluster.paused?.toString() ?? '—'}</Pill></Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Control Plane endpoint</Table.Th>
                <Table.Td>{cluster.controlPlaneEndpoint?.host ?? '—'}:{cluster.controlPlaneEndpoint?.port ?? '—'}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Pod Network</Table.Th>
                <Table.Td>{cluster.clusterNetwork?.pods?.cidrBlocks?.join(', ') ?? '—'}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Service Network</Table.Th>
                <Table.Td>{cluster.clusterNetwork?.services?.cidrBlocks?.join(', ') ?? '—'}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Service ExternalIPs</Table.Th>
                <Table.Td>{cluster.clusterNetwork?.services?.externalIPs?.join(', ') ?? '—'}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Service NodePort Range</Table.Th>
                <Table.Td>{cluster.clusterNetwork?.services?.nodePortRange ?? '—'}</Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
        <Space h="md" />
        <Panel title="Machine Deployments" content={
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
                cluster.topology?.machineDeployments?.map((md, i) => (
                  <Table.Tr className={sourceSans400.className} key={md.name ?? `md-${i}`}>
                    <Table.Td>{md.name ?? '—'}</Table.Td>
                    <Table.Td>{md.class ?? '—'}</Table.Td>
                    <Table.Td>{md.replicas ?? '—'}</Table.Td>
                    <Table.Td>{md.strategy?.type ?? '—'}</Table.Td>
                  </Table.Tr>
                ))
              }
            </Table.Tbody>
          </Table>
        } />
      </GridCol>
      <GridCol span={{base: 12, md: 6}}>
        <Panel title="Cluster Class" content={
          <Table variant="vertical">
            <Table.Tbody className="text-sm">
              <Table.Tr>
                <Table.Th w={230}>Kubernetes Version</Table.Th>
                <Table.Td>{cluster.topology?.kubernetesVersion ?? '—'}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Cluster Class Name</Table.Th>
                <Table.Td>{cluster.topology?.className ?? '—'}</Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
        <Space h="md" />
        <ConditionsTable conditions={cluster.status?.conditions ?? []}/>
        </GridCol>
      </Grid>
  )
}