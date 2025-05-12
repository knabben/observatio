import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {roboto} from "@/fonts";
import Panel from "@/app/ui/dashboard/utils/panel";
import React from "react";
import {Card, Chip, Grid, GridCol, Indicator, Pill, SimpleGrid, Space, Table} from "@mantine/core";
import { XMarkIcon } from '@heroicons/react/24/outline';

export default function MachineInfraDetails({
  machine
}: {machine: MachineInfraType}) {
  return (
    <GridCol className={roboto.className} span={12}>
      <Card withBorder shadow="sm" padding="lg" radius="md">
        <SimpleGrid className="text-center" cols={2}>
          <div>
            <span className="font-bold">Name: </span>
            {
              machine.ready
                ? <Indicator offset={-3} inline withBorder position="top-end" color="green" size={10}> {machine.name} </Indicator>
                : <Indicator  offset={-3} inline withBorder position="top-end" color="red" size={10}> {machine.name} </Indicator>
            }
          </div>
          <div><span className="font-bold">Age:</span> {machine.created}</div>
        </SimpleGrid>
      </Card>
      <Space h="md" />
      <Grid>
        <GridCol span={6}>
          <Panel title="Specification" content={
            <Table
              variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Namespace</Table.Th>
                  <Table.Td>{machine.namespace}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Provider</Table.Th>
                  <Table.Td>{machine.providerID}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>failureDomain</Table.Th>
                  <Table.Td>{machine.failureDomain}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>PowerOffMode</Table.Th>
                  <Table.Td>{machine.powerOffMode}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Template</Table.Th>
                  <Table.Td>{machine.template}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Clone Mode</Table.Th>
                  <Table.Td>{machine.cloneMode}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Age</Table.Th>
                  <Table.Td>{machine.created}</Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
          } />
        </GridCol>
        <GridCol span={6}>
          <Panel title="Machine details" content={
            <Table
              variant="vertical">
              <Table.Tbody className="text-sm">
                <Table.Tr>
                  <Table.Th w={260}>Num CPUs</Table.Th>
                  <Table.Td>{machine.numCPUs}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>CPU Per Socket</Table.Th>
                  <Table.Td>{machine.numCoresPerSocket}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Memory</Table.Th>
                  <Table.Td>{machine.memoryMiB}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Disk GiB</Table.Th>
                  <Table.Td>{machine.diskGiB}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Failure Reason</Table.Th>
                  <Table.Td>{machine.failureReason}</Table.Td>
                </Table.Tr>
                <Table.Tr>
                  <Table.Th>Failure Message</Table.Th>
                  <Table.Td>{machine.failureMessage}</Table.Td>
                </Table.Tr>
              </Table.Tbody>
            </Table>
          } />
          <Space h="md" />
          <Panel title="Machine conditions" content={
          <Table variant="vertical">
            <Table.Tbody className="text-sm">
              {
                machine.conditions?.map((condition, ic) => (
                  <Table.Tr key={ic}>
                    <Table.Td>
                      {
                        condition.status
                        ? <Chip key={ic} className="p-1" defaultChecked color="teal" variant="light">{condition.type}</Chip>
                        : <Chip key={ic} defaultChecked icon={<XMarkIcon />} color="red" variant="light">{condition.type}</Chip>
                      }
                    </Table.Td>
                    <Table.Td>{condition.lastTransitionTime}</Table.Td>
                  </Table.Tr>
                ))
              }
            </Table.Tbody>
          </Table>
          } />
        </GridCol>
      </Grid>
    </GridCol>
  )
}
