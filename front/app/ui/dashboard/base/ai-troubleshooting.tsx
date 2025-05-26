'use client';

import Panel from "@/app/ui/dashboard/utils/panel";
import {
  Chip,
  Table,
  Text,
  Notification,
  SimpleGrid,
  Button,
  Textarea,
  Space,
  Stack,
  Grid,
  GridCol
} from "@mantine/core";
import {IconX} from "@tabler/icons-react";
import React, {useEffect, useState} from "react";
import {Conditions} from "@/app/ui/dashboard/base/types";
import {postAIAnalysis} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";

type AIResponse = {
  description: string;
  solution: string;
}

export default function AITroubleshooting({
  conditions,
}: {conditions: Conditions[]}) {
  const [reasons, setReason] = useState("")
  const [loading, setLoading] = useState(false)
  const [aiResponse, setAiResponse] = useState<AIResponse>({description: "", solution: ""})

  useEffect(() => {
    const uniqueReasons = new Set(
      conditions
        ?.filter(condition => condition.reason)
        .map( (condition) => {
          const mapper = condition.reason+" of type "+condition.type
          if (condition.message != undefined) {
            return mapper + " with message: " + condition.message
          }
          return mapper
        })
    );
    setReason(Array.from(uniqueReasons).join(', '));
  }, [conditions]);

  async function requestIA() {
    try {
      setLoading(true);
      const response = await postAIAnalysis(reasons)
      setAiResponse(response);
      setLoading(false);
    } catch (error) {
      console.error('Error analyzing machine:', error);
    }
  }
  
  return (
    <Grid justify="flex-start" align="flex-start">

      <GridCol span={6}>
        <Notification withCloseButton={false} title="Status & Troubleshooting assistant" color="#a1f54d">

      <div>
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
        </div>
      </Notification>
      </GridCol>
      <GridCol span={6}>
        {
          <div>
            <Stack align="flex-end" className="text-center">
              <Button bg="#a1f54d" c="#000" variant="filled" onClick={requestIA}>Get Help with AI Agent!</Button>
            </Stack>
          </div>
        }
      </GridCol>
    </Grid>
  )
}