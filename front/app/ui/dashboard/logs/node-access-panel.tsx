'use client';

import React, {useEffect, useState} from 'react';
import {Code, Group, Stack, Text} from '@mantine/core';
import {IconAlertTriangle} from '@tabler/icons-react';
import {getNodeAccess, NodeAccessInfo} from '@/app/lib/data';
import {ObjectRef} from '@/app/ui/dashboard/shared/use-day2-ops';

interface NodeAccessPanelProps {
  objectRef: ObjectRef;
}

/**
 * Static SSH connection instructions for a Machine's node (FR-021, FR-022): a command and a
 * disclaimer only — never a live terminal, never a credential input of any kind. Observātiō does
 * not store or manage SSH credentials.
 */
export function NodeAccessPanel({objectRef}: NodeAccessPanelProps) {
  const [info, setInfo] = useState<NodeAccessInfo | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    getNodeAccess(objectRef)
      .then((result) => {
        if (!cancelled) setInfo(result);
      })
      .catch(() => {
        if (!cancelled) setError('Node access details could not be determined for this machine.');
      });
    return () => {
      cancelled = true;
    };
  }, [objectRef]);

  if (error) {
    return (
      <Group gap="xs">
        <IconAlertTriangle size={16} color="var(--mantine-color-red-6)"/>
        <Text size="sm" c="red.6">{error}</Text>
      </Group>
    );
  }

  if (!info) return <Text size="sm" c="dimmed">Loading…</Text>;

  return (
    <Stack gap={4}>
      <Code>{info.command}</Code>
      <Text size="xs" c="dimmed">{info.note}</Text>
    </Stack>
  );
}
