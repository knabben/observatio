'use client';

import React from 'react';
import {Table, Card, Text, Badge, Group} from '@mantine/core';
import {getClusterClasses} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";
import {roboto} from "@/fonts";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import { XMarkIcon } from '@heroicons/react/24/outline';
import {Conditions} from "@/app/ui/dashboard/base/types";
import {EmptyState} from "@/app/ui/dashboard/shared/empty-state";
import {useFetchState} from "@/app/ui/dashboard/shared/use-fetch-state";

type ClusterClass = {
  name?: string,
  namespace?: string,
  generation?: bigint,
  conditions?: Conditions[]
}

const TABLE_HEADERS = ['Name', 'Namespace', 'Updates', 'Status'] as const;

export const useClusterClasses = () => {
  const {data: clusterClasses, isLoading, error} = useFetchState<ClusterClass[]>(
    getClusterClasses,
    [],
    'Failed to load cluster classes',
  );
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
    <Table.Td>{clusterClass.name ?? '—'}</Table.Td>
    <Table.Td>{clusterClass.namespace ?? '—'}</Table.Td>
    <Table.Td>{clusterClass.generation?.toString() ?? '—'}</Table.Td>
    <Table.Td rowSpan={1}>
      {/* Read-only status display — a Chip is an interactive toggle and must never represent status */}
      <Group gap="xs">
        {(clusterClass.conditions ?? []).map((condition, index) => (
          condition.status?.toLowerCase() === 'true'
            ? <Badge key={index} color="teal" variant="light">{condition.type}</Badge>
            : (
              <Badge key={index} color="red" variant="light" leftSection={<XMarkIcon width={12}/>}>
                {condition.type}
              </Badge>
            )
        ))}
      </Group>
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
      {!isLoading && !error && clusterClasses.length === 0 && <EmptyState label="No cluster classes found"/>}
      {!isLoading && !error && clusterClasses.length > 0 && (
        <Table striped highlightOnHover>
          <TableHeader/>
          <Table.Tbody className="text-sm">
            {clusterClasses.map((clusterClass, i) => (
              <ClusterClassRow
                key={clusterClass.name ? `${clusterClass.namespace}-${clusterClass.name}` : `row-${i}`}
                clusterClass={clusterClass}
              />
            ))}
          </Table.Tbody>
        </Table>
      )}
    </Card>
  );
}