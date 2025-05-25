'use client';

import React from "react";

import {Table, Indicator, Badge} from '@mantine/core';
import { GridCol } from '@mantine/core';

import {roboto} from '@/fonts';
import {ClusterInfraType} from '@/app/ui/dashboard/components/clusters/types';

export default function ClusterInfraTable({
  clusters, select
}: {
  clusters: ClusterInfraType[]
  select: (cluster: ClusterInfraType) => void
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
            <Table.Th>Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            clusters?.map( (cluster: ClusterInfraType, i) => (
              <Table.Tr className={roboto.className} key={i}>
                <Table.Td>
                  <a className="cursor-pointer hover:opacity-70" onClick={() => select(cluster)}>{cluster.metadata?.name}</a>
                </Table.Td>
                <Table.Td>
                  <Badge variant="light" color="gray"> {cluster.metadata?.namespace} </Badge>
                </Table.Td>
                <Table.Td>{cluster.cluster}</Table.Td>
                <Table.Td>{cluster.server}</Table.Td>
                <Table.Td>{cluster.age}</Table.Td>
                <Table.Td ta="center">
                  {cluster.status?.ready
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