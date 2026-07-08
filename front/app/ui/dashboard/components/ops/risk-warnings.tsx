'use client';

import React from 'react';
import {Alert, Stack, Text, Title} from '@mantine/core';
import {IconAlertTriangle, IconHelpCircle} from '@tabler/icons-react';
import {RiskKind, RiskWarning} from '@/app/ui/dashboard/shared/use-day2-ops';

const RISK_KIND_LABELS: Record<RiskKind, string> = {
  cert_expiry: 'Certificate expiry',
  stalled_rollout: 'Stalled rollout',
  version_skew: 'Provider/CRD version skew',
  drift: 'Infrastructure drift',
};

interface RiskWarningsProps {
  risks: RiskWarning[];
}

/**
 * Proactively detected risk warnings (US3), grouped by kind, shown alongside the health rollups.
 * A risk with checkStatus "not_evaluable" is shown as an explicit "could not be checked" notice
 * rather than being silently omitted (FR-018).
 */
export function RiskWarnings({risks}: RiskWarningsProps) {
  if (risks.length === 0) return null;

  const grouped = new Map<RiskKind, RiskWarning[]>();
  for (const risk of risks) {
    const list = grouped.get(risk.kind) ?? [];
    list.push(risk);
    grouped.set(risk.kind, list);
  }

  return (
    <Stack gap="sm" mt="md">
      <Title order={5} ta="center">Proactive risk warnings</Title>
      {Array.from(grouped.entries()).map(([kind, warnings]) => (
        <Alert key={kind} color="orange" icon={<IconAlertTriangle size={16}/>} title={RISK_KIND_LABELS[kind]}>
          <Stack gap={4}>
            {warnings.map((risk, index) => (
              <div key={`${risk.objectRef.namespace}/${risk.objectRef.name}-${index}`}>
                {risk.checkStatus === 'not_evaluable' ? (
                  <Text size="xs" c="dimmed">
                    <IconHelpCircle size={12} style={{verticalAlign: 'text-bottom'}}/>{' '}
                    {risk.objectRef.name}: check could not be performed
                  </Text>
                ) : (
                  <Text size="xs">
                    {risk.detail}
                    {risk.likelyCause && ` — likely cause: ${risk.likelyCause}`}
                  </Text>
                )}
              </div>
            ))}
          </Stack>
        </Alert>
      ))}
    </Stack>
  );
}
