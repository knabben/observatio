import React from "react";
import {Card} from "@mantine/core";
import {openSans, roboto} from "@/fonts";
import Header from "@/app/ui/dashboard/utils/header";

export default function Panel({ title, content }: { title: string, content: any }) {
  return (
    <>
      <Card className={roboto.className} shadow="sm" padding="lg" radius="md" withBorder>
        <Header title={title} />
        {content}
      </Card>
    </>
  )
}
