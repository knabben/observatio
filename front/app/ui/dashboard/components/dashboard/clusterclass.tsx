'use client';

import React, {useState, useEffect} from 'react';
import {Table, Card, Text, Chip} from '@mantine/core';
import {getClusterClasses} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {roboto} from "@/fonts";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import { XMarkIcon } from '@heroicons/react/24/outline';
import {Conditions} from "@/app/ui/dashboard/base/types";

type ClusterClass = {
  name: string,
  namespace: string,
  generation: bigint,
  conditions: Conditions[]
}

const TABLE_HEADERS = ['Name', 'Namespace', 'Updates', 'Status'] as const;

export const useClusterClasses = () => {
  const [clusterClasses, setClusterClasses] = useState<ClusterClass[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleFetchError = (error: Error) => {
    console.error('Failed to fetch cluster summary:', error);
    setError('Failed to load cluster summary');
    setIsLoading(false);
  };

  useEffect(() => {
    const fetchClusterClasses = async () => {
      try {
        setIsLoading(true);
        const response = await getClusterClasses();
        setClusterClasses(response);
      } catch (error) {
        handleFetchError(error as Error)
      } finally {
        setIsLoading(false);
      }
    };
    fetchClusterClasses();
  }, []);

  return {clusterClasses, isLoading, error};
};

const TableHeader: React.FC = () => (
  <Table.Thead className="text-sm">
    <Table.Tr>
      {TABLE_HEADERS.map((header) => (
        <Table.Th key={header}>{header}</Table.Th>
      ))}
    </Table.Tr>
  </Table.Thead>
);

const ClusterClassRow: React.FC<{ clusterClass: ClusterClass }> = ({clusterClass}) => (
  <Table.Tr>
    <Table.Td>{clusterClass.name}</Table.Td>
    <Table.Td>{clusterClass.namespace}</Table.Td>
    <Table.Td>{clusterClass.generation}</Table.Td>
    <Table.Td rowSpan={1}>
      {
        clusterClass.conditions.map((condition, index) => (
          condition.status.toLowerCase() === 'true'
          ? <Chip key={index} className="p-1" defaultChecked color="teal" variant="light">{condition.type}</Chip>
          : <Chip key={index} defaultChecked icon={<XMarkIcon />} color="red" variant="light">{condition.type}</Chip>
        ))
      }
    </Table.Td>
  </Table.Tr>
);

/**
 * Functional component that renders a list of cluster classes.
 * It fetches cluster classes data and handles three states:
 * loading, error, and successfully fetched data.
 * Depending on the state, it displays a loader, error message,
 * or a table of cluster classes.
 */
export default function ClusterClassLister() {
  const {clusterClasses, isLoading, error} = useClusterClasses();

  return (
    <Card shadow="md" className={roboto.className} radius="md" withBorder>
      <Header title="Cluster Class"/>
      {isLoading && <CenteredLoader />}
      {error && <Text c="red">{error}</Text>}
      {!isLoading && !error && (
        <Table striped highlightOnHover>
          <TableHeader/>
          <Table.Tbody className="text-sm">
            {clusterClasses.map((clusterClass) => (
              <ClusterClassRow
                key={`${clusterClass.namespace}-${clusterClass.name}`}
                clusterClass={clusterClass}
              />
            ))}
          </Table.Tbody>
        </Table>
      )}
    </Card>
  );
}