'use client';

import React, {useState, useEffect} from 'react';
import {getComponentsVersion} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {roboto} from "@/fonts";
import {Card, Table, Text} from '@mantine/core';
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";

type Component = {
    name: string,
    kind: string,
    version: string,
}

const TABLE_HEADERS = [
  {key: 'name', label: 'Name'},
  {key: 'kind', label: 'Kind'},
  {key: 'version', label: 'Versions'}
] as const;

const useComponentVersions = () => {
  const [component, setComponent] = useState<Component[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleFetchError = (error: Error) => {
    console.error('Failed to fetch cluster summary:', error);
    setError('Failed to load cluster summary');
    setIsLoading(false);
  };

  useEffect(() => {
    const fetchVersions= async () => {
      try {
        setIsLoading(true);
        const response = await getComponentsVersion();
        setComponent(response)
      } catch (error) {
        handleFetchError(error as Error)
      } finally {
        setIsLoading(false);
      }
    };
    fetchVersions();
  }, []);

  return {component, isLoading, error};
};

const VersionsTable = ({components}: { components: Component[] }) => (
  <Table striped highlightOnHover>
    <Table.Thead className="text-sm">
      <Table.Tr>
        {TABLE_HEADERS.map(header => (
          <Table.Th key={header.key}>{header.label}</Table.Th>
        ))}
      </Table.Tr>
    </Table.Thead>
    <Table.Tbody className="text-sm">
      {components.map((component) => (
        <Table.Tr key={component.name}>
          <Table.Td>{component.name}</Table.Td>
          <Table.Td>{component.kind}</Table.Td>
          <Table.Td>{component.version}</Table.Td>
        </Table.Tr>
      ))}
    </Table.Tbody>
  </Table>
);

/**
 * This component fetches and displays a list of components along with their names, kinds, and versions
 * in a tabular format. The data is fetched asynchronously and updates when the component mounts.
 */
export default function ClusterVersions() {
  const {component, isLoading, error} = useComponentVersions();

  return (
    <Card shadow="md" className={roboto.className} radius="md" withBorder>
      <Header title="Component Versions"/>
      {isLoading && <CenteredLoader />}
      {error && <Text c="red">{error}</Text>}
      {!isLoading && !error && <VersionsTable components={component}/>}
    </Card>
  );
}
