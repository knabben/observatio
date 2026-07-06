'use client';

import React from 'react';
import {Badge} from '@mantine/core';
import {useInfraCapability} from '@/app/ui/dashboard/shared/infra-capability-context';

const LABELS: Record<'docker' | 'vsphere', {label: string; color: string}> = {
  docker: {label: 'Docker', color: 'blue'},
  vsphere: {label: 'vSphere', color: 'grape'},
};

/**
 * Per-row infrastructure provider indicator, showing the provider name plus its detected
 * version (e.g. "Docker v1.10.10"), or a neutral "Unknown" badge for an unrecognized or
 * absent provider — never blank, and never a guess.
 */
export function ProviderBadge({provider}: {provider?: string}) {
  const capability = useInfraCapability();

  if (provider !== 'docker' && provider !== 'vsphere') {
    return <Badge variant="light" color="gray">Unknown</Badge>;
  }

  const {label, color} = LABELS[provider];
  const version = capability[provider]?.version;
  return (
    <Badge variant="light" color={color}>
      {version ? `${label} ${version}` : label}
    </Badge>
  );
}
