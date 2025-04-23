'use client';

import { Input } from '@mantine/core';

export default function Search({ placeholder }: { placeholder: string }) {
  return (
    <div>
      <Input placeholder={placeholder} />
    </div>
  );
}
