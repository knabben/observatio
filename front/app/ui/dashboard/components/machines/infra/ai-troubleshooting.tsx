'use client';


import Panel from "@/app/ui/dashboard/utils/panel";
import {Chip, Table, Text, Notification} from "@mantine/core";
import {IconX} from "@tabler/icons-react";
import React from "react";
import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";

export default function AITroubleshooting({
  machine,
}: {machine: MachineInfraType}) {
  return (
    <Notification withCloseButton={false} title="AI troubleshooting assistant" color="#275479">
      <Panel title="Machine conditions" content={
        <Table variant="vertical">
          <Table.Tbody className="text-sm">
            {
              machine.status.conditions?.map((condition, ic) => (
                <Table.Tr key={ic}>
                  <Table.Td>
                    {
                      condition.status.toLowerCase() == "true"
                        ? <Chip key={ic} className="p-1" color="teal" variant="light">{condition.type}</Chip>
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
                    <Text c="#d40805"> {condition.message}</Text>
                  </Table.Td>
                </Table.Tr>
              ))
            }
          </Table.Tbody>
        </Table>
      } />
    </Notification>
  )
}