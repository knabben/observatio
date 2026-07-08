'use client';

import React from 'react';
import {Alert, Group, Text} from '@mantine/core';
import {IconAlertOctagon, IconAlertTriangle, IconInfoCircle} from '@tabler/icons-react';
import {FailureSeverity, SeverityLevel} from '@/app/ui/dashboard/shared/use-day2-ops';

const SEVERITY_ORDER: SeverityLevel[] = ['management_critical', 'provider_degraded', 'needs_investigation', 'self_healing'];

interface SeverityBannerProps {
  severities: FailureSeverity[];
}

/**
 * Top-level, hard-to-miss banner for the highest-urgency detected FailureSeverity (FR-015).
 * Self-healing activity is rendered informationally — no alert role, no red/orange styling — so
 * it is never mistaken for something requiring action (FR-012, SC-005).
 */
export function SeverityBanner({severities}: SeverityBannerProps) {
  if (severities.length === 0) return null;

  const highest = SEVERITY_ORDER
    .map((level) => severities.find((s) => s.level === level))
    .find((s): s is FailureSeverity => s != null);

  if (!highest) return null;

  if (highest.level === 'self_healing') {
    return (
      <Group gap="xs" justify="center" mb="md">
        <IconInfoCircle size={16} color="var(--mantine-color-blue-6)"/>
        <Text size="sm" c="dimmed">{highest.reason}</Text>
      </Group>
    );
  }

  const isCritical = highest.level === 'management_critical' || highest.level === 'provider_degraded';
  return (
    <Alert
      role="alert"
      color={isCritical ? 'red' : 'orange'}
      icon={isCritical ? <IconAlertOctagon size={18}/> : <IconAlertTriangle size={18}/>}
      title={isCritical ? 'Critical' : 'Needs investigation'}
      mb="md"
    >
      {highest.reason}
    </Alert>
  );
}
