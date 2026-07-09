'use client';

import React from "react";

import {Pill, SimpleGrid, Space, Table} from "@mantine/core";
import Panel from "@/app/ui/dashboard/utils/panel";
import ConditionsTable from "@/app/ui/dashboard/shared/conditions-table";
import {ClusterClassType} from "@/app/ui/dashboard/components/clusterclasses/types";

export default function Specification({
  cc,
 }: {cc: ClusterClassType}) {
  return (
    <SimpleGrid cols={{base: 1, md: 2}}>
      <div style={{gridColumn: '1 / -1'}}>
        <Panel title="Specification" content={
          <Table
            variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Namespace</Table.Th>
                  <Table.Td>{cc.namespace ?? '—'}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Generation</Table.Th>
                  <Table.Td><Pill size="sm">{cc.generation ?? '—'}</Pill></Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
        } />
      </div>
      <div style={{gridColumn: '1 / -1'}}>
        <Space h="md" />
        <ConditionsTable conditions={cc.conditions ?? []}/>
      </div>
    </SimpleGrid>
  )
}
