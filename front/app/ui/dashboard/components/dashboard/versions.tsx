'use client';

import React, {useState, useEffect} from 'react';
import { Card, Table, Text, Divider } from '@mantine/core';
import {getComponentsVersion} from "@/app/lib/data";
import {sourceCodePro400} from "@/fonts";

type Component = {
    name: string,
    kind: string,
    version: string,
}

// ClusterVersions: Cluster components versions and enumeration.
export default function ClusterVersions() {
  const [clusterVersions, setClusterVersions] = useState<Component[]>([]);
  useEffect( () => {
    const fetch = async  () => {
      setClusterVersions(await getComponentsVersion())
    }
    fetch().catch( (e) => { console.error('error', e) })
  }, [])

  return (
    <Card shadow="md"  radius="md" withBorder>
      <Text size="lg" tt="uppercase" className={sourceCodePro400.className} c="#8feb83" ta="center">
        Component Versions
      </Text>
      <Divider my="sm" variant="dashed" />
      <Table striped highlightOnHover withTableBorder withColumnBorders>
        <Table.Thead>
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>Kind</Table.Th>
          <Table.Th>Versions</Table.Th>
        </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {
            clusterVersions.map((component) => (
              <Table.Tr key={component.name}>
                <Table.Td>{component.name}</Table.Td>
                <Table.Td>{component.kind}</Table.Td>
                <Table.Td>{component.version}</Table.Td>
              </Table.Tr>
            ))
          }
        </Table.Tbody>
      </Table>
    </Card>
  );
}
