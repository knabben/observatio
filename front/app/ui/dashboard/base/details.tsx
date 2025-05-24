import {Card, GridCol, Space, Tabs, TabsPanel} from "@mantine/core";
import React from "react";
import {roboto} from "@/fonts";

/**
 * Example usage:
 *
 * const tabs = [
 *   {
 *     label: "Details",
 *     content: <DetailContent data={someData} />
 *   },
 *   {
 *     label: "History",
 *     content: <HistoryContent data={someData} />
 *   }
 * ];
 *
 * const headerRender = (obj: MyType) => (
 *   <div>
 *     <h1>{obj.name}</h1>
 *     <p>{obj.description}</p>
 *   </div>
 * );
 *
 * <ObjectDetails
 *   object={myObject}
 *   headerRenderer={headerRender}
 *   tabs={tabs}
 * />
 */
export default function ObjectDetails<T>({
  object,
  headerRenderer,
  tabs,
}: {
  object: T,
  headerRenderer: (object: T) => React.ReactNode,
  tabs: { label: string, content: (object: T) => React.ReactNode }[]
}) {
  return (
    <GridCol className={roboto.className} span={12}>
      <Card withBorder shadow="sm" padding="lg" radius="md">
        {headerRenderer(object)}
      </Card>
      <Space h="md" />
      <Tabs mb="md" color="#48654a" defaultValue={tabs[0]?.label}>
        <Tabs.List>
          {tabs.map((tab) => (
            <Tabs.Tab key={tab.label} value={tab.label}>{tab.label}</Tabs.Tab>
          ))}
        </Tabs.List>
        {tabs.map((tab) => (
          <TabsPanel key={tab.label} value={tab.label}>
            <Space h="lg"/>
            {tab.content(object)}
          </TabsPanel>
        ))}
      </Tabs>
    </GridCol>
  )
}
