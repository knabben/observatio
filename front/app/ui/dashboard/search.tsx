'use client';

import React from 'react';

import { Input } from '@mantine/core';
import { GridCol } from '@mantine/core';

export default function Search({
  placeholder,
  onChange,
  value
}: { placeholder: string, onChange?: (e:React.ChangeEvent<HTMLInputElement>) => void, value: string }) {
  return (
    <GridCol span={4}>
      <Input
        value={value}
        onChange={onChange}
        placeholder={placeholder}
      />
    </GridCol>
  );
}
