'use client';

import React, {useState, useEffect} from 'react';
import {Table, Card, Text} from '@mantine/core';
import {getClusterClasses} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {roboto} from "@/fonts";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";

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
    {clusterClass.conditions.map((condition) => (
      <Table.Td key={`${condition.type}-${condition.lastTransitionTime}`} rowSpan={1}>
        {condition.type} - {condition.status}
      </Table.Td>
    ))}
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