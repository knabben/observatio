'use client';

import {useEffect} from "react";
import {Button, Center, Stack, Text, Title} from "@mantine/core";

/**
 * Root App-Router error boundary. Any render throw in the tree lands here instead of a
 * blank white screen, and the operator can attempt recovery.
 */
export default function Error({error, reset}: {error: Error & {digest?: string}; reset: () => void}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <Center mih="100vh" p="xl">
      <Stack align="center" gap="md" maw={520}>
        <Title order={2}>Something went wrong</Title>
        <Text c="dimmed" ta="center">
          The dashboard hit an unexpected error while rendering this view. Your cluster is unaffected.
        </Text>
        <Button onClick={reset} color="teal">Try again</Button>
      </Stack>
    </Center>
  );
}
