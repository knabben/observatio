'use client';

import Panel from "@/app/ui/dashboard/utils/panel";
import {
  Avatar,
  ActionIcon,
  Text,
  Group,
  Button,
  Card,
  ScrollArea,
  Textarea,
  Stack,
  Box,
  Paper,
  Grid,
  GridCol,
  Table,
  Chip
} from '@mantine/core';
import {IconX, IconArrowsDiagonal, IconArrowsMinimize} from "@tabler/icons-react";
import React, {useEffect, useRef, useState} from "react";
import {v4 as uuidv4} from "uuid";
import {ReadyState} from "react-use-websocket";
import {Conditions} from "@/app/ui/dashboard/base/types";
import {WebSocket, WS_URL_CHATBOT} from "@/app/lib/websocket";

/** Bounded time to await an AI response before resetting the loading indicator. */
const AI_RESPONSE_TIMEOUT_MS = 30_000;

export default function AITroubleshooting({
  objectType,
  objectName,
  objectNamespace,
  conditions,
}: {
  objectType: string,
  objectName: string,
  objectNamespace: string,
  conditions: Conditions[]
}) {
  const [request, setRequest] = useState("")
  const [expanded, setExpanded] = useState(false)

  useEffect(() => {
    const broken = new Set(
      conditions?.filter(condition => condition.reason && condition.status != "True")
        .map( (condition) => {
          const mapper = "On " + objectName + " of type " + objectType + ", running on namespace " + objectNamespace +
            " it is failing with condition " + condition.reason
          if (condition.message != undefined) {
            return mapper + " with message: " + condition.message
          }
          return mapper
        })
    );
    if (broken.size > 0) {
      setRequest(Array.from(broken).join(', '));
    }
  }, [conditions, objectType, objectName, objectNamespace]);

  return (
    <Grid justify="flex-start" align="flex-start">
      <GridCol span={expanded ? 12 : {base: 12, md: 6}}>
        <ChatBot expanded={expanded} setExpanded={setExpanded} request={request}/>
      </GridCol>
      <GridCol span={{base: 12, md: 6}}>
        <Panel title="Object conditions" content={
          <Table variant="vertical">
            <Table.Tbody className="text-sm">
              {
                conditions?.map((condition, ic) => (
                  <Table.Tr key={ic}>
                    <Table.Td>
                      {
                        condition.status?.toLowerCase() === "true"
                          ? <Chip key={ic} defaultChecked className="p-1" color="teal" variant="light">{condition.type}</Chip>
                          : (
                            <div>
                              {[condition.type, condition.reason].map((text, index) => (
                                <Chip
                                  key={`${ic}-${index}`}
                                  icon={<IconX size={16}/>}
                                  className="p-1"
                                  color="red"
                                  defaultChecked={true}
                                  variant="light"
                                >
                                  {text}
                                </Chip>
                              ))}
                            </div>
                          )
                      }
                    </Table.Td>
                    <Table.Td>{condition.lastTransitionTime}</Table.Td>
                    <Table.Td>{condition.severity}</Table.Td>
                    <Table.Td className="break-all text-xs">
                      {
                        condition.status?.toLowerCase() === "true"
                          ? <Text size="sm" fw={700} className="text-xs break-all">{condition.message}</Text>
                          : <Text size="sm" c="red">{condition.message}</Text>
                      }
                    </Table.Td>
                  </Table.Tr>
                ))
              }
            </Table.Tbody>
          </Table>
        } />
      </GridCol>

    </Grid>
  )
}

type WSRequest = {
  id: string,
  type: string,
  agent_id: string,
  content: string,
  timestamp: string,
  actor: string,
}

function ChatBot({
  request,
  expanded,
  setExpanded,
}: {
  request: string,
  expanded: boolean,
  setExpanded: React.Dispatch<React.SetStateAction<boolean>>
}) {
  const [messages, setMessages] = useState<WSRequest[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [aiRequest, setAIRequest] = useState(request)
  const scrollRef = useRef<HTMLDivElement | null>(null)
  const responseTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket(WS_URL_CHATBOT)

  useEffect(() => {
    setAIRequest(request)
  }, [request]);

  useEffect(() => {
    if (lastJsonMessage) {
      const response = lastJsonMessage as WSRequest;
      setMessages(prevMessages => [...prevMessages, response]);
      setIsLoading(false)
      if (responseTimeoutRef.current) clearTimeout(responseTimeoutRef.current);
    }
  }, [lastJsonMessage])

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTo({top: scrollRef.current.scrollHeight, behavior: 'smooth'});
    }
  }, [messages])

  // Clear any pending response timeout on unmount so it never fires against a stale component.
  useEffect(() => () => {
    if (responseTimeoutRef.current) clearTimeout(responseTimeoutRef.current);
  }, []);

  async function requestIA() {
    if (aiRequest === "" || readyState !== ReadyState.OPEN) {
      return;
    }
    try {
      setAIRequest('')
      const request: WSRequest = {
        id: uuidv4(),
        type: "chatbot",
        content: aiRequest,
        actor: "user",
        agent_id: "cluster-agent",
        timestamp: new Date().toLocaleDateString('en-US', {
          month: '2-digit',
          day: '2-digit',
          year: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
          hour12: false
        })
      }
      setIsLoading(true)
      sendJsonMessage(request)
      setMessages(prevMessages => [...prevMessages, request])
      responseTimeoutRef.current = setTimeout(() => setIsLoading(false), AI_RESPONSE_TIMEOUT_MS);
    } catch (error) {
      console.error('Error analyzing machine:', error);
      setIsLoading(false)
    }
  }

  return (
    <Card
      h="100%"
      radius="lg"
      withBorder
      style={{
        background: 'linear-gradient(135deg, #0f0f23 0%, #1a1a3e 100%)',
        border: '1px solid var(--mantine-color-brand-8)',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <Group justify="space-between" p="md" style={{borderBottom: '1px solid var(--mantine-color-brand-8)'}}>
        <Text fw={700} c="var(--mantine-color-brand-4)">AI Troubleshooting</Text>
        <ActionIcon
          onClick={() => setExpanded((prev) => !prev)}
          color="green"
          aria-label={expanded ? 'Collapse AI troubleshooting panel' : 'Expand AI troubleshooting panel'}
        >
          {expanded
            ? <IconArrowsMinimize style={{ width: '70%', height: '70%' }} stroke={1.5} />
            : <IconArrowsDiagonal style={{ width: '70%', height: '70%' }} stroke={1.5} />}
        </ActionIcon>
      </Group>
      <ScrollArea flex={1} p="md" viewportRef={(ref) => {
        scrollRef.current = ref
      }}>
        <Stack gap="md">
          {
            messages.map((message) => (
              <Group
                key={message.id}
                align="flex-start"
                justify={message.actor === 'user' ? 'flex-end' : 'flex-start'}
                gap="sm"
              >
              {message.actor === 'agent' && (
                <Avatar size="sm" color="rgba(0, 153, 204, 0.5)" variant="outline" radius="xl">BOT</Avatar>
              )}
              <Paper
                p="md"
                radius="lg"
                maw="90%"
                style={{
                  background: message.actor === 'user'
                    ? 'rgba(0, 212, 170, 0.1)'
                    : 'rgba(0, 153, 204, 0.2)',
                  border: message.actor === 'user'
                    ? '1px solid rgba(0, 212, 170, 0.3)'
                    : '1px solid rgba(0, 153, 204, 0.4)',
                }}
              >
                {/* Plain-text render only — the content is untrusted (user/AI supplied) and
                    React escapes it automatically; no dangerouslySetInnerHTML. */}
                <Text size="sm" className="break-all" style={{ lineHeight: 1.5, whiteSpace: 'pre-wrap' }}>
                  {message.content}
                </Text>
                <Text size="xs" c="dimmed" mt="xs">
                  {message.timestamp}
                </Text>
              </Paper>
              {message.actor === 'user' && (
                <Avatar size="sm" variant="outline" color="rgba(0, 212, 170, 0.5)" radius="xl">USR</Avatar>
              )}
          </Group>
            ))
          }
        </Stack>
      </ScrollArea>
      <Box className="content-right" p="md" style={{ borderTop: '1px solid #4a4a6a' }}>
        <Group>
          <Textarea
            flex={1}
            placeholder="Ask about this issue or request specific actions..."
            className="min-w-full"
            value={aiRequest}
            onChange={(e) => setAIRequest(e.target.value)}
            radius="xl"
            styles={{
              input: {
                height: '100px',
                border: '1px solid #48654a',
                color: '#e0e0e0',
                '&:focus': {
                  borderColor: '#00d4aa'
                }
              }
            }}
          />
          <Button onClick={requestIA} disabled={isLoading || readyState !== ReadyState.OPEN} bg="var(--mantine-color-brand-4)" c="#000" variant="filled">Send!</Button>
        </Group>
      </Box>
    </Card>
  );
};
