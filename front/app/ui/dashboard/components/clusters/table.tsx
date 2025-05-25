'use client';

import React from "react";

import {Table, Indicator, Badge} from '@mantine/core';
import { GridCol } from '@mantine/core';

import {roboto} from '@/fonts';
import {ClusterType} from '@/app/ui/dashboard/components/clusters/types';

export default function ClusterTable({
  clusters, select
}: {
  clusters: ClusterType[],
  select: (cluster: ClusterType) => void
}) {
  return (
    <GridCol span={12}>
      <Table highlightOnHover>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Namespace</Table.Th>
            <Table.Th>Version</Table.Th>
            <Table.Th>Phase</Table.Th>
            <Table.Th>Age</Table.Th>
            <Table.Th>Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            clusters?.map( (cluster: ClusterType, i) => (
              <Table.Tr className={roboto.className} key={i}>
                <Table.Td>
                  <a className="cursor-pointer hover:opacity-70" onClick={() => select(cluster)}>{cluster.metadata.name}</a>
                </Table.Td>
                <Table.Td>
                  <Badge variant="light" color="gray"> {cluster.metadata.namespace} </Badge>
                </Table.Td>
                <Table.Td>{cluster.topology?.kubernetesVersion}</Table.Td>
                <Table.Td>{cluster.status.phase}</Table.Td>
                <Table.Td ta="center">{cluster.age}</Table.Td>
                <Table.Td ta="center">
                  {cluster.status.controlPlaneReady && cluster.status.infrastructureReady
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