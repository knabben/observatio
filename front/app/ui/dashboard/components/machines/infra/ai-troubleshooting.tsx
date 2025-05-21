'use client';

import Panel from "@/app/ui/dashboard/utils/panel";
import {Chip, Table, Text, Notification, SimpleGrid, Button, Textarea, Space, Stack} from "@mantine/core";
import {IconX} from "@tabler/icons-react";
import React, {useEffect, useState} from "react";
import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {postAIAnalysis} from "@/app/lib/data";
import Header from "@/app/ui/dashboard/utils/header";

type AIResponse = {
  description: string;
  solution: string;
}

export default function AITroubleshooting({
  machine,
}: {machine: MachineInfraType}) {
  const [reasons, setReason] = useState("")
  const [loading, setLoading] = useState(false)
  const [aiResponse, setAiResponse] = useState<AIResponse>({description: "", solution: ""})

  useEffect(() => {
    const uniqueReasons = new Set(
      machine.status.conditions
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
  }, [machine.status.conditions]);

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
    <Notification withCloseButton={false} title="Status & Troubleshooting assistant" color="#a1f54d">
      <Space h="lg" />
      <SimpleGrid cols={2}>
        {
          reasons &&
          <div>
            <Stack align="flex-end" className="text-center">
              <Textarea
                className="min-w-full"
                styles={{input: {height: '150px'}}}
                value={reasons}
                onChange={(e) => setReason(e.target.value)} />
              {loading
                ? <Text className="text-center text-white">Analyzing...</Text>
                : <Button bg="#a1f54d" c="#000" variant="filled" onClick={requestIA}>Get Help!</Button>
              }
            </Stack>
            {aiResponse.description && aiResponse.solution && (
              <>
                <Space h="md"/>
                <Notification color="gray" withCloseButton={false}>
                  <Header title="Analysis and Description" />
                  <div className="text-white" dangerouslySetInnerHTML={{__html: aiResponse.description}}/>
                </Notification>
                <Space h="md"/>
                <Notification withCloseButton={false} color="#304a47">
                  <Header title="How to fix" />
                  <div className="text-white" dangerouslySetInnerHTML={{__html: aiResponse.solution}}/>
                </Notification>
              </>
            )}
          </div>
        }
      <div>
      <Panel title="Machine conditions" content={
        <Table variant="vertical">
          <Table.Tbody className="text-sm">
            {
              machine.status.conditions?.map((condition, ic) => (
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
                    <Text c="red" className="text-bold"> {condition.message}</Text>
                  </Table.Td>
                </Table.Tr>
              ))
            }
          </Table.Tbody>
        </Table>
      } />
        </div>
      </SimpleGrid>
    </Notification>
  )
}