'use client';

import React from "react";

import {Pill, SimpleGrid, Space, Table} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import ConditionsTable from "@/app/ui/dashboard/shared/conditions-table";
import {KubeadmControlPlaneType} from "@/app/ui/dashboard/components/kubeadmcontrolplanes/types";

export default function Specification({
  kcp,
 }: {kcp: KubeadmControlPlaneType}) {
  return (
    <SimpleGrid cols={{base: 1, md: 2}}>
      <div>
        <Panel title="Specification" content={
          <Table
            variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Namespace</Table.Th>
                  <Table.Td>{kcp.metadata?.namespace ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Cluster</Table.Th>
                  <Table.Td>{kcp.cluster ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Version</Table.Th>
                  <Table.Td>{kcp.version ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Desired Replicas</Table.Th>
                  <Table.Td><Pill size="sm">{kcp.replicas ?? '—'}</Pill></Table.Td>
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
                <Table.Td><Pill size="sm">{kcp.status?.readyReplicas ?? '—'}</Pill></Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Updated Replicas</Table.Th>
                <Table.Td><Pill size="sm">{kcp.status?.updatedReplicas ?? '—'}</Pill></Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Unavailable Replicas</Table.Th>
                <Table.Td><Pill size="sm">{kcp.status?.unavailableReplicas ?? '—'}</Pill></Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Initialized</Table.Th>
                <Table.Td>{kcp.status?.initialized == null ? '—' : String(kcp.status.initialized)}</Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
      </div>
      <div style={{gridColumn: '1 / -1'}}>
        <Space h="md" />
        <ConditionsTable conditions={kcp.status?.conditions ?? []}/>
      </div>
    </SimpleGrid>
  )
}
