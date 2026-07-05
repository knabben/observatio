'use client';

import {useEffect, useState} from 'react';

export interface FetchState<T> {
  data: T;
  isLoading: boolean;
  error: string | null;
}

/**
 * Dedupes the fetch/isLoading/error boilerplate previously copy-pasted (with a copy-pasted
 * error string) across every dashboard widget. `errorMessage` is shown verbatim on failure,
 * so each caller supplies its own — no more "Failed to fetch cluster summary" on a components
 * or topology fetch failure.
 */
export function useFetchState<T>(fetchFn: () => Promise<T>, initial: T, errorMessage: string): FetchState<T> {
  const [data, setData] = useState<T>(initial);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        setIsLoading(true);
        const response = await fetchFn();
        if (!cancelled) setData(response);
      } catch (err) {
        console.error(`${errorMessage}:`, err);
        if (!cancelled) setError(errorMessage);
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return {data, isLoading, error};
}
