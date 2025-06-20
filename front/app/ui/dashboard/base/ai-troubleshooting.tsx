'use client';

import Panel from "@/app/ui/dashboard/utils/panel";
import {
  AppShell,
  Avatar,
  ActionIcon,
  Text,
  Group,
  Button,
  Card,
  Notification,
  ScrollArea,
  Textarea,
  Stack,
  Box,
  Paper,
  Container,
  Grid,
  GridCol,
  Table,
  Chip
} from '@mantine/core';
import {IconX, IconArrowsDiagonal} from "@tabler/icons-react";
import React, {useEffect, useState} from "react";
import {Conditions} from "@/app/ui/dashboard/base/types";
import {WebSocket, WS_URL_CHATBOT} from "@/app/lib/websocket";

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
  const [expandChat, setExpandChat] = useState(6)

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
      <GridCol span={expandChat}>

        <ChatBot setExpand={setExpandChat} request={request}/>
      </GridCol>
      <GridCol span={6}>
        <Panel title="Object conditions" content={
          <Table variant="vertical">
            <Table.Tbody className="text-sm">
              {
                conditions?.map((condition, ic) => (
                  <Table.Tr key={ic}>
                    <Table.Td>
                      {
                        condition.status.toLowerCase() == "true"
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
                        condition.status.toLowerCase() == "true"
                          ? <Text size="sm" className="text-xs text-bold break-all">{condition.message}</Text>
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
  setExpand,
}: {
  request: string,
  setExpand: React.Dispatch<React.SetStateAction<number>>
}) {
  const [messages, setMessages] = useState<WSRequest[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [aiRequest, setAIRequest] = useState(request)
  const scrollRef = React.useRef<HTMLDivElement | null>(null)
  const {sendJsonMessage, lastJsonMessage} = WebSocket(WS_URL_CHATBOT)

  useEffect(() => {
    setAIRequest(request)
  }, [request]);

  useEffect(() => {
    if (lastJsonMessage) {
      const response = lastJsonMessage as WSRequest;
      setMessages(prevMessages => [...prevMessages, response]);
      setIsLoading(false)
    }
  }, [lastJsonMessage])

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTo({top: scrollRef.current.scrollHeight, behavior: 'smooth'});
    }
  }, [messages])

  async function requestIA() {
    try {
      if (aiRequest != "") {
        setAIRequest('')
        const request: WSRequest = {
          id: crypto.randomUUID(),
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
        setMessages([...messages, request])
      }
    } catch (error) {
      console.error('Error analyzing machine:', error);
      setIsLoading(false)
    }
  }

  return (
    <AppShell
      header={{ height: 60 }}
      padding={0}
      styles={{
        main: {
          background: 'linear-gradient(135deg, #0f0f23 0%, #1a1a3e 100%)',
          minHeight: '100vh'
        }
      }}
    >
      <Notification withCloseButton={false} title="AI Troubleshooting" color="#a1f54d">

        <Container fluid p="md" h="calc(90vh - 60px)">
          <Card
            h="100%"
            radius="lg"
            style={{
              border: '1px solid #48654a',
              display: 'flex',
              flexDirection: 'column'
            }}
          >
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
                    <Text  size="sm" style={{ lineHeight: 1.5 }}>
                      <div className="break-all wrap-break-word" dangerouslySetInnerHTML={{__html: message.content}}/>
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
              <ActionIcon onClick={() => setExpand(12)} color="green" aria-label="Settings">
                <IconArrowsDiagonal style={{ width: '70%', height: '70%' }} stroke={1.5} />
              </ActionIcon>
              <Button onClick={requestIA} disabled={isLoading} bg="#a1f54d" c="#000" variant="filled">Send!</Button>
            </Group>
            </Box>
          </Card>
        </Container>
      </Notification>
    </AppShell>
  );
};
