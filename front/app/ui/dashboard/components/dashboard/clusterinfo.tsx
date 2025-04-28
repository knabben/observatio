'use client';

import React, {useState, useEffect} from 'react';
import { Table, Card } from '@mantine/core';
import {getClusterInformation} from "@/app/lib/data";
import {roboto} from "@/fonts";
import Header from "@/app/ui/dashboard/utils/header";

type service = {
  name: string,
  path: string
}

// ClusterInfo: Cluster details and access URLs.
export default function ClusterInfo() {
  const [clusterInfo, setClusterInfo] = useState<service[]>([])
  useEffect( () => {
    const fetch = async  () => {
      setClusterInfo(await getClusterInformation())
    }
    fetch().catch( (e) => { console.error('error', e) })
  }, [])

  return (
    <Card shadow="md" className={roboto.className} radius="md" withBorder>
      <Header title="Cluster Information" />
      <Table striped highlightOnHover>
        <Table.Thead className="text-sm">
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>URL</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
          {
            clusterInfo.map((service) => (
              <Table.Tr key={service.name}>
                <Table.Td>{service.name}</Table.Td>
                <Table.Td><a href={service.path}>{service.path}</a></Table.Td>
              </Table.Tr>
            ))
          }
        </Table.Tbody>
      </Table>
    </Card>
  );
}