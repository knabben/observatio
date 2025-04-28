'use client';

import React, {useState, useEffect} from 'react';
import { Table, Card, Text, Divider } from '@mantine/core';
import {getClusterClasses} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {roboto, sourceCodePro400} from "@/fonts";

type Conditions = {
  type: string,
  status: boolean,
  lastTransitionTime: string,
}

type ClusterClass = {
  name: string,
  namespace: string,
  generation: bigint,
  conditions: Conditions[]
}

// ClusterClass: details the cluster classes existent in the cluster.
export default function ClusterClass() {
  const [clusterClass, setClusterClass] = useState<ClusterClass[]>([])
  useEffect( () => {
    const fetch = async  () => {
      setClusterClass(await getClusterClasses())
    }
    fetch().catch( (e) => { console.error('error', e) })
  }, [])

  return (
    <Card shadow="md" className={roboto.className}  radius="md" withBorder>
    <Header title="Cluster Class" />
    <Table striped highlightOnHover>
      <Table.Thead className="text-sm">
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>Namespace</Table.Th>
          <Table.Th>Updates</Table.Th>
          <Table.Th>Status</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody className="text-sm">
        {
          clusterClass.map( (cc, i) => (
            <Table.Tr key={i}>
              <Table.Td>{cc.name}</Table.Td>
              <Table.Td>{cc.namespace}</Table.Td>
              <Table.Td>{cc.generation}</Table.Td>
              {cc.conditions.map((condition, i) => (
                  <Table.Td key={i} rowSpan={1}>
                    {condition.type} - {condition.status}
                  </Table.Td>
              ))}
            </Table.Tr>
          ))
        }
      </Table.Tbody>
    </Table>
    </Card>
  );
}