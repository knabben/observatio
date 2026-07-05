import React from "react";
import {Center, Stack, Text, ThemeIcon} from "@mantine/core";
import {IconInbox} from "@tabler/icons-react";

interface EmptyStateProps {
  label: string;
  hint?: string;
}

/**
 * Labeled empty-state message shown when a collection has zero items — replaces the
 * former header-only tables and empty charts.
 */
export const EmptyState: React.FC<EmptyStateProps> = ({label, hint}) => (
  <Center mih={160} w="100%">
    <Stack align="center" gap="xs">
      <ThemeIcon variant="light" color="gray" size="xl" radius="xl">
        <IconInbox size={22}/>
      </ThemeIcon>
      <Text c="dimmed" fw={500}>{label}</Text>
      {hint && <Text c="dimmed" size="sm">{hint}</Text>}
    </Stack>
  </Center>
);
