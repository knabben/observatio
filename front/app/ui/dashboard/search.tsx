'use client';

import {useState} from 'react';
import { Input } from '@mantine/core';

export default function Search({
  placeholder,
  onChange,
  value
}: { placeholder: string, onChange?: (e: any) => void, value: string }) {
  return (
    <div>
      <Input
        value={value}
        onChange={onChange}
        placeholder={placeholder}
      />
    </div>
  );
}
