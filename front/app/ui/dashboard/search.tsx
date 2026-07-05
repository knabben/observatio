'use client';

import React from 'react';
import {TextInput} from '@mantine/core';
import {IconSearch} from '@tabler/icons-react';

/**
 * Live, client-side text filter over the current list (filters by name as you type).
 * Previously labeled "Search" but implemented as a pick-one dropdown that never
 * filtered anything — this now does what the label says.
 */
export default function Search({
  value,
  onChange,
  placeholder = 'Search by name…',
}: {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}) {
  return (
    <TextInput
      value={value}
      onChange={(e) => onChange(e.currentTarget.value)}
      placeholder={placeholder}
      leftSection={<IconSearch size={16} aria-hidden="true"/>}
      aria-label="Search by name"
    />
  );
}
