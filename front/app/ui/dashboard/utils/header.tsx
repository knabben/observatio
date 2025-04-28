import React from "react";
import {Divider, Title} from "@mantine/core";
import {openSans} from "@/fonts";

export default function Header({ title }: { title: string }) {
  return (
    <>
      <Title c="#8feb83" ta="center" className={openSans.className} order={4}>{title}</Title>
      <Divider my="sm" variant="dashed" />
    </>
  )
}
