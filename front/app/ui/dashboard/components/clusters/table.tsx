'use client';

import React from "react";
import { Table } from '@mantine/core';
import { GridCol } from '@mantine/core';

import {ClusterType} from '@/app/ui/dashboard/components/clusters/types'

export default function ClusterTable({
  clusters
}: {
  clusters: ClusterType[]
}) {
  return (
    <GridCol span={12}>
      <Table striped highlightOnHover withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>ClusterClass</Table.Th>
            <Table.Th>Conditions</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {
            clusters?.map( (cluster) => (
              <Table.Tr key={cluster.name}>
                <Table.Td>{cluster.name}</Table.Td>
                <Table.Td>{cluster.conditions.length}</Table.Td>
              </Table.Tr>
            ))
          }
        </Table.Tbody>
      </Table>
    </GridCol>
  )
}