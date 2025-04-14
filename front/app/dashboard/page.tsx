'use client';

import { sourceCodePro400 } from '@/fonts'
import { Card, Table, Grid, Text, Title} from '@mantine/core';

export default async function Home() {
  const data = await fetch("http://localhost:8080/api/clusters/info")
  const services = await data.json()

  return (
    <main>
      <Grid grow>
        <Grid.Col span={12}>
          <Title order={2} tt="capitalize">
            Dashboard
          </Title>
        </Grid.Col>
      </Grid>
      <Grid grow>
        <Grid.Col span={4}>
          <Card shadow="sm" padding="lg" radius="md" withBorder>
            <Text>
              CAPI Versions
            </Text>
          </Card>
        </Grid.Col>
        <Grid.Col span={8}>
          <Card shadow="sm" padding="lg" radius="md" withBorder>
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
          </Card>
        </Grid.Col>
      </Grid>
    </main>
  );
}
