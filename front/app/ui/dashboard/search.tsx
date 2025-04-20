'use client';

import { Input } from '@mantine/core';
import { useSearchParams, usePathname, useRouter } from 'next/navigation';

export default function Search({ placeholder }: { placeholder: string }) {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  function handleSearch(term: string) {
    const params = new URLSearchParams(searchParams);
    if (term) {
      params.set('query', term);
    } else {
      params.delete('query');
    }
    replace(`${pathname}?${params.toString()}`);
  }

  return (
    <div>
      <Input
        onChange={(e) => handleSearch(e.currentTarget.value)}
        placeholder={placeholder}
        defaultValue={searchParams.get('query')?.toString()}
      />
    </div>
  );
}
