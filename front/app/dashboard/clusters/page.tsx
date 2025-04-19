import { getClusterList } from "@/app/lib/data";
import { Grid, GridCol, Title, Space } from '@mantine/core';
import ClusterLister from '@/app/ui/dashboard/components/ClusterLister'

export default async function Clusters() {
  const clusters = await getClusterList()
  return (
    <div>
      <main>
        <Title order={3}>Clusters</Title>
        <Space h="md" />
        <Grid grow>
          <GridCol span={12}>
            <ClusterLister clusterList={clusters} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}