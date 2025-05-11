'use client';

import React from "react";

import { Table, Indicator } from '@mantine/core';
import { GridCol } from '@mantine/core';

import {roboto} from '@/fonts';
import {ClusterInfraType} from '@/app/ui/dashboard/components/clusters/types';

export default function ClusterInfraTable({
  clusters
}: {
  clusters: ClusterInfraType[]
}) {
  return (
    <GridCol span={12}>
      <Table highlightOnHover>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Namespace</Table.Th>
            <Table.Th>Cluster</Table.Th>
            <Table.Th>Server</Table.Th>
            <Table.Th>Age</Table.Th>
            <Table.Th ta="center">Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            clusters?.map( (cluster: ClusterInfraType) => (
              <Table.Tr className={roboto.className} key={cluster.name}>
                <Table.Td>{cluster.name}</Table.Td>
                <Table.Td>{cluster.namespace}</Table.Td>
                <Table.Td>{cluster.cluster}</Table.Td>
                <Table.Td>{cluster.server}</Table.Td>
                <Table.Td>{cluster.created}</Table.Td>
                <Table.Td ta="center">
                  {cluster.ready
                    ? <Indicator inline processing color="green" size={15}/>
                    : <Indicator inline processing color="red" size={15}/>
                  }</Table.Td>
              </Table.Tr>
            ))
          }
        </Table.Tbody>
      </Table>
    </GridCol>
  )
}