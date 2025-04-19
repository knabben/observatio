import {getMachines} from "@/app/lib/data";
import { Title, Grid, GridCol, Space } from '@mantine/core';
import MachineLister from "@/app/ui/dashboard/components/MachineLister";

export default async function Machines() {
  const machines = await getMachines()
  return (
    <div>
      <main>
        <Title order={3}>Machine Deployments</Title>
        <Space h="md" />
        <Grid grow>
          <GridCol span={12}>
            <MachineLister machines={machines} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}