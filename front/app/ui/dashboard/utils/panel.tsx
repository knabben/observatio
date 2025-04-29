/* eslint-disable @typescript-eslint/no-explicit-any */
import React from "react";
import {Card} from "@mantine/core";
import {roboto} from "@/fonts";
import Header from "@/app/ui/dashboard/utils/header";

export default function Panel({ title, content }: { title: string, content: any }) {
  return (
    <>
      <Card className={roboto.className} padding="md" radius="sm" withBorder>
        <Header title={title} />
        {content}
      </Card>
    </>
  )
}
