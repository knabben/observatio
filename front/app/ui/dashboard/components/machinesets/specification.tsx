'use client';

import React from "react";

import {Pill, SimpleGrid, Space, Table} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import ConditionsTable from "@/app/ui/dashboard/shared/conditions-table";
import {MachineSetType} from "@/app/ui/dashboard/components/machinesets/types";

export default function Specification({
  ms,
 }: {ms: MachineSetType}) {
  return (
    <SimpleGrid cols={{base: 1, md: 2}}>
      <div>
        <Panel title="Specification" content={
          <Table
            variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Namespace</Table.Th>
                  <Table.Td>{ms.metadata?.namespace ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Cluster</Table.Th>
                  <Table.Td>{ms.cluster ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Machine Deployment</Table.Th>
                  <Table.Td>{ms.machineDeployment ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Desired Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{ms.replicas ?? '—'}</Pill></Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
        } />
      </div>
      <div>
        <Panel title="Status" content={
          <Table
            variant="vertical">
            <Table.Tbody className="text-sm">
              <Table.Tr>
                <Table.Th w={260}>Ready Replicas</Table.Th>
                <Table.Td><Pill size="sm">{ms.status?.readyReplicas ?? '—'}</Pill></Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Available Replicas</Table.Th>
                <Table.Td><Pill size="sm">{ms.status?.availableReplicas ?? '—'}</Pill></Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Fully Labeled Replicas</Table.Th>
                <Table.Td><Pill size="sm">{ms.status?.fullyLabeledReplicas ?? '—'}</Pill></Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
      </div>
      <div style={{gridColumn: '1 / -1'}}>
        <Space h="md" />
        <ConditionsTable conditions={ms.status?.conditions ?? []}/>
      </div>
    </SimpleGrid>
  )
}
