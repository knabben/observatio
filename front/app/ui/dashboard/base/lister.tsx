/* eslint-disable @typescript-eslint/no-explicit-any */
'use client';

import React, {useEffect} from 'react';
import { sourceCodePro400 } from "@/fonts";

import {useState} from "react";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";

import { Grid, GridCol, Title } from '@mantine/core';
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {IconArrowBigLeft} from "@tabler/icons-react";
import {FilterItems} from "@/app/dashboard/utils";

interface BaseListerProps<T extends object> {
  objectType: string;
  items: T[];
  renderDetails: (item: T) => React.ReactNode;
  renderTable: (items: T[], handleSelect: any) => React.ReactNode;
  title: string;
  titleLink?: string;
}

type GenericMeta = {
  metadata: {
    name: string
  }
}

export default function BaseLister<T extends object>({
  objectType,
  items: initialItems,
  renderDetails,
  renderTable,
  title,
}: BaseListerProps<T>) {
  const [items, setItems] = useState<T[]>(initialItems)
  const [selected, setSelected] = useState('')
  const [loading, setLoading] = useState(true)

  const handleSelect = (item: T | null) => {
    if (item === null || Object.keys(item as object).length === 0) {
      setSelected('')
      return
    }

    const itemAsAny = item as any
    setSelected(itemAsAny?.metadata?.name || '')
  }

  const filteredItem: T | undefined = selected
    ? FilterItems(selected, items)
    : undefined;

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, objectType, sendJsonMessage)
  }, [readyState, sendJsonMessage, objectType])

  useEffect(() => {
    const newItems = receiveAndPopulate(lastJsonMessage, [...items]).sort(
      (a: GenericMeta, b: GenericMeta) => a?.metadata?.name.localeCompare(b?.metadata?.name)
    )
    setItems(newItems)
    setLoading(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastJsonMessage])

  if (loading) {
    return <CenteredLoader/>;
  }

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol span={9}>
        <Title className={sourceCodePro400.className} order={2}>
          {title}
        </Title>
      </GridCol>
      <GridCol span={3} className="flex justify-end items-center">
        { selected &&
          <div>
            <IconArrowBigLeft onClick={() => handleSelect(null)} size={32} className="cursor-pointer hover:opacity-70"/>
          </div>
        }
      </GridCol>
      {
        filteredItem
          ? renderDetails(filteredItem)
          : renderTable(items, handleSelect)
      }
    </Grid>
  );
} 