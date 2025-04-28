'use client';

import React, {useState, useEffect} from 'react';
import {Card, Table} from '@mantine/core';
import {getComponentsVersion} from "@/app/lib/data";
import {roboto} from "@/fonts";
import Header from "@/app/ui/dashboard/utils/header";

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
    <Card shadow="md" className={roboto.className}  radius="md" withBorder>
      <Header title="Component Versions" />
      <Table striped highlightOnHover>
        <Table.Thead className="text-sm">
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>Kind</Table.Th>
          <Table.Th>Versions</Table.Th>
        </Table.Tr>
        </Table.Thead>
        <Table.Tbody className="text-sm">
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
