import { Table } from '@mantine/core';

async function getServices() {
  const res = await fetch("http://localhost:8080/api/clusters/info")
  return res.json()
}

export default async function ClusterInfo() {
  const response = getServices()
  const [services] = await Promise.all([response])

  return (
    <Table striped highlightOnHover withTableBorder withColumnBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>Name</Table.Th>
          <Table.Th>URL</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {
          services["services"].map( (service: { name: string, path: string }, i) => (
            <Table.Tr key={service.name}>
              <Table.Td>{service.name}</Table.Td>
              <Table.Td><a href={service.path}>{service.path}</a></Table.Td>
            </Table.Tr>
          ))
        }
      </Table.Tbody>
    </Table>
  );
}
