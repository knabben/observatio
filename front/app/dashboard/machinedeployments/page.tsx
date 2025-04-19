import {getMachinesDeployments} from "@/app/lib/data";
import { Title, Grid, GridCol, Space } from '@mantine/core';
import MachineDeploymentLister from "@/app/ui/dashboard/components/MachineDeploymentLister";

export default async function MachineDeployments() {
  const mds = await getMachinesDeployments()
  return (
    <div>
      <main>
        <Title order={3}>Machine Deployments</Title>
        <Space h="md" />
        <Grid grow>
          <GridCol span={12}>
            <MachineDeploymentLister mds={mds} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}