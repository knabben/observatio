import {ClusterInfraDockerType} from "@/app/ui/dashboard/components/clusters/types";
import Panel from "@/app/ui/dashboard/utils/panel";
import {Table} from "@mantine/core";
import React from "react";

export default function DockerSpecification({
  cluster,
}: { cluster: ClusterInfraDockerType }) {
  return (
    <Panel title="Specification" content={
      <Table variant="vertical">
        <Table.Tbody className="text-sm">
          <Table.Tr>
            <Table.Th w={260}>Cluster</Table.Th>
            <Table.Td>{cluster.cluster ?? '—'}</Table.Td>
          </Table.Tr>
          <Table.Tr>
            <Table.Th w={260}>Load Balancer IP</Table.Th>
            <Table.Td>{cluster.loadBalancerIP ?? '—'}</Table.Td>
          </Table.Tr>
          <Table.Tr>
            <Table.Th w={260}>Ready</Table.Th>
            <Table.Td>{cluster.ready ? 'true' : 'false'}</Table.Td>
          </Table.Tr>
        </Table.Tbody>
      </Table>
    }/>
  )
}
