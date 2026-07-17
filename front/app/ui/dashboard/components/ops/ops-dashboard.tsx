'use client';

import React, {useMemo, useState} from 'react';
import {Alert, Button, Grid, GridCol, Group, Stack, Title} from '@mantine/core';
import {IconAlertTriangle} from '@tabler/icons-react';
import {useDay2Ops} from '@/app/ui/dashboard/shared/use-day2-ops';
import {Category, HealthRollup} from '@/app/ui/dashboard/shared/use-day2-ops';
import {HealthRollupCard} from './health-rollup-card';
import {BackupHealthCard} from './backup-health-card';
import {ToolSourcesCard} from './tool-sources-card';
import {DebuggingPath} from './debugging-path';
import {RiskWarnings} from './risk-warnings';
import {SeverityBanner} from './severity-banner';

const CATEGORY_FILTERS: {label: string; value: Category | 'all'}[] = [
  {label: 'All', value: 'all'},
  {label: 'Clusters', value: 'cluster'},
  {label: 'Machine Deployments', value: 'machine_deployment'},
  {label: 'Machines', value: 'machine'},
];

/**
 * Top-level Day-2 Ops landing view: consolidated per-category health rollups (FR-001, FR-002),
 * narrowed to a single category in place on click (FR-003), with an explicit "data unavailable"
 * banner instead of a false all-clear/blank state when the source connection is lost (FR-017).
 */
export function OpsDashboard() {
  const {data} = useDay2Ops();
  const [filter, setFilter] = useState<Category | 'all'>('all');

  const visibleRollups = useMemo<HealthRollup[]>(() => {
    if (filter === 'all') return data.rollups;
    return data.rollups.filter((r) => r.category === filter);
  }, [data.rollups, filter]);

  return (
    <>
      {data.sourceUnavailable && (
        <Alert
          color="red"
          icon={<IconAlertTriangle size={18}/>}
          title="Data unavailable"
          mb="md"
        >
          The connection to the management cluster was lost. The information below may be stale.
        </Alert>
      )}
      {/* Once sourceUnavailable is true, the banner above already covers the management-critical
          severity it derives from (research.md); showing both would be redundant (FR-015/FR-017). */}
      {!data.sourceUnavailable && <SeverityBanner severities={data.severities}/>}
      <Group justify="center" mb="md" gap="xs">
        {CATEGORY_FILTERS.map((f) => (
          <Button
            key={f.value}
            variant={filter === f.value ? 'filled' : 'default'}
            size="xs"
            onClick={() => setFilter(f.value)}
          >
            {f.label}
          </Button>
        ))}
      </Group>
      <Grid grow justify="center">
        {visibleRollups.map((rollup) => (
          <GridCol key={rollup.category} span={{base: 12, sm: filter === 'all' ? 4 : 12}}>
            <HealthRollupCard rollup={rollup}/>
          </GridCol>
        ))}
        {filter === 'all' && (
          <GridCol span={{base: 12, sm: 4}}>
            <BackupHealthCard health={data.backupHealth}/>
          </GridCol>
        )}
        {filter === 'all' && (
          <GridCol span={{base: 12, sm: 4}}>
            <ToolSourcesCard/>
          </GridCol>
        )}
      </Grid>
      {(filter === 'all' || filter === 'machine') && data.debugPaths.length > 0 && (
        <Stack gap="sm" mt="md">
          <Title order={5} ta="center">Unhealthy machines — likely cause</Title>
          {data.debugPaths.map((path) => (
            <DebuggingPath key={`${path.objectRef.namespace}/${path.objectRef.name}`} path={path}/>
          ))}
        </Stack>
      )}
      <RiskWarnings risks={data.risks}/>
    </>
  );
}
