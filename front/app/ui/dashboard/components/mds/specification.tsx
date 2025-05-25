'use client';

import React from "react";

import {Pill, SimpleGrid, Table} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";

export default function Specification({
  md,
 }: {md: MachineDeploymentType}) {
  return (
    <SimpleGrid cols={2}>
      <div>
        <Panel title="Specification" content={
          <Table
            variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Namespace</Table.Th>
                  <Table.Td>{md.metadata?.namespace}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Template version</Table.Th>
                  <Table.Td>{md.templateversion}</Table.Td>
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
                  <Table.Td><Pill size="sm">{md.status.readyReplicas}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Updated Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{md.status.updatedReplicas}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Unavailable Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{md.status.unavailableReplicas}</Pill></Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
        } />
      </div>
      <div>
        <Panel title="References" content={
          <Table
            variant="vertical">
            <Table.Tbody className="text-sm">
              <Table.Tr>
                <Table.Th w={260}>Bootstrap name</Table.Th>
                <Table.Td>{md.templateBootstrap?.configRef.name}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}>Bootstrap version</Table.Th>
                <Table.Td>{md.templateBootstrap?.configRef?.apiVersion}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th w={260}>Bootstrap kind</Table.Th>
                <Table.Td>{md.templateBootstrap?.configRef?.kind}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Infrastructure name</Table.Th>
                <Table.Td>{md.templateInfrastructureRef?.name}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Infrastructure version</Table.Th>
                <Table.Td>{md.templateInfrastructureRef?.apiVersion}</Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Infrastructure kind</Table.Th>
                <Table.Td>{md.templateInfrastructureRef?.kind}</Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
      </div>
    </SimpleGrid>
  )
}