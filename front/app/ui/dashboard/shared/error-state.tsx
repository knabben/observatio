import React from "react";
import {Button, Center, Stack, Text, ThemeIcon} from "@mantine/core";
import {IconAlertTriangle, IconRefresh} from "@tabler/icons-react";

interface ErrorStateProps {
  message: string;
  onRetry?: () => void;
}

/**
 * Actionable error message with an optional keyboard-accessible retry control, shown
 * when a data channel fails (HTTP error, dropped/exhausted WebSocket) instead of a
 * silent failure or a perpetual loader.
 */
export const ErrorState: React.FC<ErrorStateProps> = ({message, onRetry}) => (
  <Center mih={160} w="100%">
    <Stack align="center" gap="xs">
      <ThemeIcon variant="light" color="red" size="xl" radius="xl">
        <IconAlertTriangle size={22}/>
      </ThemeIcon>
      <Text c="red.7" fw={500} ta="center">{message}</Text>
      {onRetry && (
        <Button
          variant="light"
          color="red"
          size="xs"
          leftSection={<IconRefresh size={14}/>}
          onClick={onRetry}
        >
          Retry
        </Button>
      )}
    </Stack>
  </Center>
);
