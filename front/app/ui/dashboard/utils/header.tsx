import React from "react";
import {Divider, Title} from "@mantine/core";
import {openSans} from "@/fonts";

export default function Header({ title }: { title: string }) {
  return (
    <>
      {/* brand.8 (#48654a) on white ≈ 6.5:1 contrast (WCAG AA); the prior #8feb83 was ≈1.5:1 and unreadable */}
      <Title c="var(--mantine-color-brand-8)" ta="center" className={openSans.className} order={4}>{title}</Title>
      <Divider my="sm" variant="dashed" />
    </>
  )
}
