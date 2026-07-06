'use client';

import {Button} from '@mantine/core';
import {IconMessageChatbot} from '@tabler/icons-react';
import React from 'react';
import {formatObjectContext, ObjectContext, useAIPanel} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

/**
 * Per-object-screen quick-action (FR-016): opens the global AI panel already pre-filled with
 * this screen's object context in one click, rather than requiring the operator to open the
 * panel separately and notice it auto-filled.
 */
export function AskAIButton({context}: { context: ObjectContext }) {
  const {open} = useAIPanel();

  return (
    <Button
      onClick={() => open(formatObjectContext(context))}
      variant="light"
      size="xs"
      leftSection={<IconMessageChatbot size={16}/>}
    >
      Ask AI about this
    </Button>
  );
}
