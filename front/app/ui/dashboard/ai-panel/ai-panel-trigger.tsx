'use client';

import {ActionIcon, Tooltip} from '@mantine/core';
import {IconMessageChatbot} from '@tabler/icons-react';
import React from 'react';
import {useAIPanel} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

/**
 * Persistent, fixed-position control reachable from any dashboard screen (FR-001) — opens the
 * app-wide AI panel pre-filled from whatever object is currently in view, or empty/general if
 * none is (FR-007).
 */
export default function AIPanelTrigger() {
  const {open, isOpen} = useAIPanel();

  if (isOpen) return null;

  return (
    <Tooltip label="Ask AI" position="left">
      <ActionIcon
        onClick={() => open()}
        size="xl"
        radius="xl"
        variant="filled"
        bg="var(--mantine-color-brand-4)"
        c="#000"
        aria-label="Open AI troubleshooting panel"
        style={{
          position: 'fixed',
          bottom: 24,
          right: 24,
          zIndex: 200,
          boxShadow: 'var(--mantine-shadow-md)',
        }}
      >
        <IconMessageChatbot size={26} stroke={1.5}/>
      </ActionIcon>
    </Tooltip>
  );
}
