'use client';

import { Table } from '@mantine/core';

// Versions: Cluster components versions and enumeration.
export default function Versions() {
    return (
    <Table striped highlightOnHover withTableBorder withColumnBorders>
        <Table.Thead>
        <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Versions</Table.Th>
        </Table.Tr>
        </Table.Thead>
    </Table>
    );
}
