'use client';

import Panel from "@/app/ui/dashboard/utils/panel";
import {
  AppShell,
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
import {IconX} from "@tabler/icons-react";
import React, {useEffect, useState} from "react";
import {Conditions} from "@/app/ui/dashboard/base/types";
import {postAIAnalysis} from "@/app/lib/data";

type AIResponse = {
  description: string;
  solution: string;
}

export default function AITroubleshooting({
  conditions,
}: {
  conditions: Conditions[]
}) {
  return (
    <Grid justify="flex-start" align="flex-start">
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
                    <Table.Td>
                      {
                        condition.status.toLowerCase() == "true"
                          ? <Text className="text-bold">{condition.message}</Text>
                          : <Text c="red" className="text-bold">{condition.message}</Text>
                      }
                    </Table.Td>
                  </Table.Tr>
                ))
              }
            </Table.Tbody>
          </Table>
        } />
      </GridCol>
      <GridCol span={6}>
        <ChatBot conditions={conditions}/>
      </GridCol>
    </Grid>
  )
}

function ChatBot({
  conditions,
}: {
  conditions: Conditions[]
}) {
  const [aiRequest, setAIRequest] = useState('')
  const [aiResponse, setAiResponse] = useState<AIResponse>({description: "", solution: ""})

  useEffect(() => {
    const broken = new Set(
      conditions?.filter(condition => condition.reason && condition.status != "True")
        .map( (condition) => {
          const mapper = condition.reason + " of type " + condition.type
          if (condition.message != undefined) {
            return mapper + " with message: " + condition.message
          }
          return mapper
        })
    );
    if (broken.size > 0) {
      setAIRequest(Array.from(broken).join(', '));
    }
  }, [conditions]);

  async function requestIA() {
    try {
      const response = await postAIAnalysis(aiRequest)
      setAiResponse(response);
    } catch (error) {
      console.error('Error analyzing machine:', error);
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
        <Container fluid p="md" h="calc(100vh - 60px)">
          <Card
            h="100%"
            radius="lg"
            style={{
              border: '1px solid #48654a',
              display: 'flex',
              flexDirection: 'column'
            }}
          >
          <ScrollArea flex={1} p="md">
            <Stack gap="md">
              { aiResponse.description != "" && aiResponse.solution != "" &&
                <>
                  <Group align="flex-start" justify='flex-start' gap="sm">
                    <Paper p="md" radius="lg" maw="80%" style={{
                      background: 'rgba(0, 212, 170, 0.1)',
                      border: '1px solid rgba(0, 212, 170, 0.3)',
                    }}>
                      <Text size="sm" style={{lineHeight: 1.5}}>
                        <Text className="text-bold" fw={700}>Description</Text>
                        <div dangerouslySetInnerHTML={{__html: aiResponse.description}}/>
                      </Text>
                      <Text size="xs" c="dimmed" mt="xs">
                        {new Date().toLocaleDateString()}
                      </Text>
                    </Paper>
                  </Group>
                  <Group align="flex-start" justify='flex-start' gap="sm">
                    <Paper p="md" radius="lg" maw="80%" style={{
                      background: 'rgba(0, 212, 170, 0.1)',
                      border: '1px solid rgba(0, 212, 170, 0.3)',
                    }}>
                      <Text size="sm" style={{lineHeight: 1.5}}>
                        <Text className="text-bold" fw={700}>Solution</Text>
                        <div dangerouslySetInnerHTML={{__html: aiResponse.solution}}/>
                      </Text>
                      <Text size="xs" c="dimmed" mt="xs">
                        {new Date().toLocaleDateString()}
                      </Text>
                    </Paper>
                  </Group>
                </>
              }
            </Stack>
          </ScrollArea>

          <Box p="md" style={{ borderTop: '1px solid #4a4a6a' }}>
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
                    height: '130px',
                    border: '1px solid #48654a',
                    color: '#e0e0e0',
                    '&:focus': {
                      borderColor: '#00d4aa'
                    }
                  }
                }}
              />
              <Button onClick={requestIA} bg="#a1f54d" c="#000" variant="filled">Send!</Button>
            </Group>
            </Box>
          </Card>
        </Container>
      </Notification>
    </AppShell>
  );
};