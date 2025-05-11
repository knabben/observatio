'use client';

import React from "react";

import { Table, Indicator } from '@mantine/core';
import { GridCol } from '@mantine/core';

import {roboto} from '@/fonts';
import {ClusterType} from '@/app/ui/dashboard/components/clusters/types';

export default function ClusterTable({
  clusters
}: {
  clusters: ClusterType[]
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
            <Table.Th ta="center">Age</Table.Th>
            <Table.Th ta="center">Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            clusters?.map( (cluster) => (
              <Table.Tr className={roboto.className} key={cluster.name}>
                <Table.Td>{cluster.name}</Table.Td>
                <Table.Td>{cluster.namespace}</Table.Td>
                <Table.Td>{cluster.clusterClass.kubernetesVersion}</Table.Td>
                <Table.Td>{cluster.phase}</Table.Td>
                <Table.Td ta="center">{cluster.created}</Table.Td>
                <Table.Td ta="center">
                  {cluster.controlPlaneReady && cluster.infrastructureReady
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