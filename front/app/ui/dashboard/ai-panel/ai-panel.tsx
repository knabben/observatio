'use client';

import {
  Avatar,
  Button,
  Drawer,
  Group,
  Loader,
  Paper,
  ScrollArea,
  Stack,
  Text,
  Textarea,
  TypographyStylesProvider,
} from '@mantine/core';
import React, {useCallback, useEffect, useRef, useState} from 'react';
import {v4 as uuidv4} from 'uuid';
import {ReadyState} from 'react-use-websocket';
import {IconAlertTriangle} from '@tabler/icons-react';
import Markdown from 'markdown-to-jsx';
import {WebSocket, WS_URL_CHATBOT} from '@/app/lib/websocket';
import {useAIPanel, WSRequest} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

/** Bounded time to await an AI response before resetting the loading indicator. */
const AI_RESPONSE_TIMEOUT_MS = 30_000;

/**
 * App-wide collapsible AI troubleshooting panel (FR-001, FR-002). Replaces the old per-object
 * embedded section: a single instance, mounted once, reachable from anywhere via a Drawer.
 * Colors use the dashboard's theme tokens throughout (FR-010) instead of the previous hardcoded
 * dark palette.
 */
export default function AIPanel() {
  const {isOpen, close, messages, setMessages, queryField, setQueryField} = useAIPanel();
  const [isLoading, setIsLoading] = useState(false);
  const [unavailable, setUnavailable] = useState(false);
  const scrollRef = useRef<HTMLDivElement | null>(null);
  const responseTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const onReconnectStop = useCallback(() => setUnavailable(true), []);
  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket(WS_URL_CHATBOT, {onReconnectStop});

  // A fresh successful connection clears any earlier "unavailable" state (e.g. the operator
  // reopened the panel after the server was reconfigured/restarted).
  useEffect(() => {
    if (readyState === ReadyState.OPEN) setUnavailable(false);
  }, [readyState]);

  useEffect(() => {
    if (!lastJsonMessage) return;
    const chunk = lastJsonMessage as WSRequest;

    // A streamed reply arrives as one or more "delta" chunks sharing the same id, followed by a
    // "done" chunk once it's finished; a chunk with no event is a one-shot complete message
    // (e.g. the "AI assistant is not available" error), handled the same way "done" clears loading.
    if (chunk.event === 'delta') {
      setMessages((prev) => {
        const last = prev[prev.length - 1];
        if (last?.id === chunk.id && last.actor === 'agent') {
          return [...prev.slice(0, -1), {...last, content: last.content + chunk.content, timestamp: chunk.timestamp}];
        }
        return [...prev, chunk];
      });
      return;
    }

    if (chunk.event !== 'done') {
      setMessages((prev) => [...prev, chunk]);
    }
    setIsLoading(false);
    if (responseTimeoutRef.current) clearTimeout(responseTimeoutRef.current);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastJsonMessage]);

  useEffect(() => {
    scrollRef.current?.scrollTo?.({top: scrollRef.current.scrollHeight, behavior: 'smooth'});
  }, [messages, isLoading]);

  // Clear any pending response timeout on unmount so it never fires against a stale component.
  useEffect(() => () => {
    if (responseTimeoutRef.current) clearTimeout(responseTimeoutRef.current);
  }, []);

  async function requestIA() {
    if (queryField === '' || readyState !== ReadyState.OPEN) {
      return;
    }
    const request = {
      id: uuidv4(),
      type: 'chatbot',
      content: queryField,
      actor: 'user',
      agent_id: 'cluster-agent',
      timestamp: new Date().toLocaleDateString('en-US', {
        month: '2-digit',
        day: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      }),
    };
    setQueryField('');
    setIsLoading(true);
    sendJsonMessage(request);
    setMessages((prev) => [...prev, request]);
    responseTimeoutRef.current = setTimeout(() => setIsLoading(false), AI_RESPONSE_TIMEOUT_MS);
  }

  return (
    <Drawer
      opened={isOpen}
      onClose={close}
      position="right"
      size="lg"
      transitionProps={{duration: 0}}
      title={<Text fw={700} c="var(--mantine-color-brand-4)">AI Troubleshooting</Text>}
      styles={{
        body: {height: 'calc(100% - 60px)', display: 'flex', flexDirection: 'column', padding: 0},
      }}
    >
      <ScrollArea flex={1} p="md" viewportRef={(ref) => {
        scrollRef.current = ref;
      }}>
        <Stack gap="md">
          {messages.map((message) => (
            <Group
              key={message.id}
              align="flex-start"
              justify={message.actor === 'user' ? 'flex-end' : 'flex-start'}
              gap="sm"
            >
              {message.actor === 'agent' && (
                <Avatar size="sm" color="gray" variant="outline" radius="xl">BOT</Avatar>
              )}
              <Paper
                p="md"
                radius="lg"
                maw="90%"
                bg={message.actor === 'user' ? 'var(--mantine-color-brand-1)' : 'var(--mantine-color-gray-1)'}
                style={{
                  border: message.actor === 'user'
                    ? '1px solid var(--mantine-color-brand-4)'
                    : '1px solid var(--mantine-color-gray-4)',
                }}
              >
                {message.actor === 'agent' ? (
                  // markdown-to-jsx parses to React elements, not dangerouslySetInnerHTML; with
                  // disableParsingRawHTML any literal "<tag>" in the untrusted, model-generated
                  // content is escaped to visible text instead of becoming a real DOM element.
                  <TypographyStylesProvider p={0} m={0} fz="sm" style={{lineHeight: 1.5, overflowWrap: 'break-word'}}>
                    <Markdown options={{disableParsingRawHTML: true}}>{message.content}</Markdown>
                  </TypographyStylesProvider>
                ) : (
                  <Text size="sm" style={{lineHeight: 1.5, whiteSpace: 'pre-wrap', overflowWrap: 'break-word'}}>
                    {message.content}
                  </Text>
                )}
                <Text size="xs" c="dimmed" mt="xs">
                  {message.timestamp}
                </Text>
              </Paper>
              {message.actor === 'user' && (
                <Avatar size="sm" variant="outline" color="var(--mantine-color-brand-6)" radius="xl">USR</Avatar>
              )}
            </Group>
          ))}
          {isLoading && (
            <Group align="flex-start" justify="flex-start" gap="sm">
              <Avatar size="sm" color="gray" variant="outline" radius="xl">BOT</Avatar>
              <Paper
                p="md"
                radius="lg"
                bg="var(--mantine-color-gray-1)"
                style={{border: '1px solid var(--mantine-color-gray-4)'}}
              >
                <Loader size="xs" type="dots" color="var(--mantine-color-brand-4)"/>
              </Paper>
            </Group>
          )}
        </Stack>
      </ScrollArea>
      <Group p="md" style={{borderTop: '1px solid var(--mantine-color-default-border)'}} wrap="nowrap" align="flex-end">
        <Textarea
          flex={1}
          placeholder="Ask about this issue or request specific actions..."
          value={queryField}
          onChange={(e) => setQueryField(e.target.value)}
          radius="md"
          autosize
          minRows={2}
          maxRows={6}
          disabled={unavailable}
        />
        <Button
          onClick={requestIA}
          disabled={isLoading || readyState !== ReadyState.OPEN || unavailable}
          bg="var(--mantine-color-brand-4)"
          c="#000"
          variant="filled"
        >
          Send
        </Button>
      </Group>
      {unavailable && (
        <Group gap="xs" px="md" pb="md" wrap="nowrap">
          <IconAlertTriangle size={16} color="var(--mantine-color-red-6)"/>
          <Text size="xs" c="red.6">
            AI assistant is not available — the server could not establish a connection (it may not
            be configured, e.g. a missing API key).
          </Text>
        </Group>
      )}
    </Drawer>
  );
}
