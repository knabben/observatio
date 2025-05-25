import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";
import Panel from "@/app/ui/dashboard/utils/panel";
import { SimpleGrid, Table } from "@mantine/core";
import React from "react";

export default function Specification({
  cluster,
}: {cluster: ClusterInfraType}) {
  return(
    <SimpleGrid cols={2}>
      <div>
        <Panel title="Specification" content={
          <Table variant="vertical">
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
      </div>
    </SimpleGrid>
  )
}