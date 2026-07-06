import React from "react";
import {Chip, Table, Text} from "@mantine/core";
import {IconX} from "@tabler/icons-react";
import Panel from "@/app/ui/dashboard/utils/panel";
import {Conditions} from "@/app/ui/dashboard/base/types";

/**
 * Object conditions table, moved here from the old per-object embedded AI Troubleshooting tab
 * (which has been replaced by the global AI panel) so status/conditions stay visible on the
 * Specification tab rather than disappearing (spec 005 FR-005).
 */
export default function ConditionsTable({conditions}: { conditions: Conditions[] }) {
  return (
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
    }/>
  );
}
