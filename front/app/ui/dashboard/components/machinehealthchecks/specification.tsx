'use client';

import React from "react";

import {Pill, SimpleGrid, Space, Table} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import ConditionsTable from "@/app/ui/dashboard/shared/conditions-table";
import {MachineHealthCheckType} from "@/app/ui/dashboard/components/machinehealthchecks/types";

export default function Specification({
  mhc,
 }: {mhc: MachineHealthCheckType}) {
  return (
    <SimpleGrid cols={{base: 1, md: 2}}>
      <div>
        <Panel title="Specification" content={
          <Table
            variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Namespace</Table.Th>
                  <Table.Td>{mhc.metadata?.namespace ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Cluster</Table.Th>
                  <Table.Td>{mhc.cluster ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Max Unhealthy</Table.Th>
                  <Table.Td>{mhc.maxUnhealthy ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Node Startup Timeout</Table.Th>
                  <Table.Td>{mhc.nodeStartupTimeout ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Expected Machines</Table.Th>
                  <Table.Td><Pill size="sm">{mhc.status?.expectedMachines ?? '—'}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Current Healthy</Table.Th>
                  <Table.Td><Pill size="sm">{mhc.status?.currentHealthy ?? '—'}</Pill></Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Remediations Allowed</Table.Th>
                  <Table.Td><Pill size="sm">{mhc.status?.remediationsAllowed ?? '—'}</Pill></Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
        } />
      </div>
      <div>
        <Panel title="Selector & Targets" content={
          <Table
            variant="vertical">
            <Table.Tbody className="text-sm">
              <Table.Tr>
                <Table.Th w={260}>Match Labels</Table.Th>
                <Table.Td>
                  {mhc.selector?.matchLabels
                    ? Object.entries(mhc.selector.matchLabels).map(([k, v]) => (
                      <Pill key={k} size="sm" mr={4}>{k}={v}</Pill>
                    ))
                    : '—'}
                </Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Unhealthy Conditions</Table.Th>
                <Table.Td>
                  {mhc.unhealthyConditions?.length
                    ? mhc.unhealthyConditions.map((c, i) => (
                      <Pill key={i} size="sm" mr={4}>{c.type}={c.status} ({c.timeout})</Pill>
                    ))
                    : '—'}
                </Table.Td>
              </Table.Tr>
              <Table.Tr>
                <Table.Th>Targets</Table.Th>
                <Table.Td>
                  {mhc.status?.targets?.length
                    ? mhc.status.targets.map((t) => <Pill key={t} size="sm" mr={4}>{t}</Pill>)
                    : '—'}
                </Table.Td>
              </Table.Tr>
            </Table.Tbody>
          </Table>
        } />
      </div>
      <div style={{gridColumn: '1 / -1'}}>
        <Space h="md" />
        <ConditionsTable conditions={mhc.status?.conditions ?? []}/>
      </div>
    </SimpleGrid>
  )
}
