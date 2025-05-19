'use client';


import Panel from "@/app/ui/dashboard/utils/panel";
import {Chip, Table, Text, Notification, SimpleGrid, Button, Textarea, Space} from "@mantine/core";
import {IconX} from "@tabler/icons-react";
import React, {useEffect, useState} from "react";
import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {postAIAnalysis} from "@/app/lib/data";

export default function AITroubleshooting({
  machine,
}: {machine: MachineInfraType}) {
  const [reasons, setReason] = useState("")
  const [aiResponse, setAiResponse] = useState<string>("")

  useEffect(() => {
    const uniqueReasons = new Set(
      machine.status.conditions
        ?.filter(condition => condition.reason)
        .map( (condition) => {
          let mapper = condition.reason+" of type "+condition.type
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
      const response = await postAIAnalysis(reasons)
      setAiResponse(response.data);
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
            <Textarea minLength={10} value={reasons} onChange={(e) => setReason(e.target.value)} readOnly/>
            <Button bg="#a1f54d" c="#000" variant="filled" onClick={requestIA}>Get Help!</Button>
            {aiResponse && (
              <>
                <Space h="md"/>
                {aiResponse}
              </>
            )}
          </div>
        }
        <div>
          <Notification withCloseButton={false} color="#b61c11">
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
                    <Text c="#b61c11" className="text-bold"> {condition.message}</Text>
                  </Table.Td>
                </Table.Tr>
              ))
            }
          </Table.Tbody>
        </Table>
      } />
          </Notification>
      </div>
</SimpleGrid>
    </Notification>
  )
}