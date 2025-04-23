'use client';

import React, { useState, useEffect } from 'react';
import { getClusterList } from "@/app/lib/data";

import Link from 'next/link';
import { Table } from '@mantine/core';
import { Grid, GridCol, Title } from '@mantine/core';

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";

type Conditions = {
  type: string,
  status: boolean,
  lastTransitionTime: string,
}

type Cluster = {
  name: string,
  hasTopology: boolean,
  conditions: Conditions[]
}

// ClusterLister: Cluster list and details component.
export default function ClusterLister() {
  const [clusters,setClusters] = useState<Cluster[]>([])
  const [query,setQuery] = useState('')
  const filteredClusters = FilterItems(query, clusters);

  useEffect(() => {
    const fetchData = async () => {
      setClusters(await getClusterList())
    }
    fetchData().catch((e) => {
      console.error('error', e)
    })
  }, [])

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/clusters">
          <Title className="hidden md:block" order={2}>
            Clusters / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <GridCol span={4}>
        <Search
          value={query}
          onChange={(e: { currentTarget: { value: React.SetStateAction<string>; }; }) => setQuery(e.currentTarget.value)}
          placeholder="Cluster name" />
      </GridCol>
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
              filteredClusters?.map( (cluster) => (
                <Table.Tr key={cluster.name}>
                  <Table.Td>{cluster.name}</Table.Td>
                  <Table.Td>{cluster.hasTopology.toString()}</Table.Td>
                  <Table.Td>{cluster.conditions.length}</Table.Td>
                </Table.Tr>
              ))
            }
          </Table.Tbody>
        </Table>
      </GridCol>
    </Grid>
  );
}