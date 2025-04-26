'use client';

import React, {useState, useEffect} from 'react';
import { Table, Card, Text, Divider } from '@mantine/core';
import {getClusterInformation} from "@/app/lib/data";
import {sourceCodePro400} from "@/fonts";

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
    <Card shadow="md"  radius="md" withBorder>
      <Text size="lg" tt="uppercase" className={sourceCodePro400.className} c="#8feb83" ta="center">
        Cluster Information
      </Text>
      <Divider my="sm" variant="dashed" />
      <Table striped highlightOnHover withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>URL</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {
            clusterInfo.map( (service) => (
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