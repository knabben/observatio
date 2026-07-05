'use client';

import React, {useState} from 'react';
import {sourceCodePro400} from "@/fonts";

import {Grid, GridCol, Title} from '@mantine/core';
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {IconArrowBigLeft} from "@tabler/icons-react";
import {FilterItems} from "@/app/dashboard/utils";
import {useResourceStream} from "@/app/ui/dashboard/shared/resource-hooks";
import {EmptyState} from "@/app/ui/dashboard/shared/empty-state";
import {ErrorState} from "@/app/ui/dashboard/shared/error-state";

interface BaseListerProps<T extends object> {
  objectType: string;
  /** Retained for API compatibility; the live stream is the source of truth. */
  items?: T[];
  renderDetails: (item: T) => React.ReactNode;
  renderTable: (items: T[], handleSelect: (item: T | null) => void) => React.ReactNode;
  title: string;
  titleLink?: string;
}

type WithMeta = {metadata?: {name?: string}};

/**
 * Live list/detail shell. Drives an explicit connecting/ready/empty/error state from
 * `useResourceStream` so a screen never hangs on a silent socket, wipes its list on an
 * empty frame, or crashes on a partial item.
 */
export default function BaseLister<T extends object>({
  objectType,
  renderDetails,
  renderTable,
  title,
}: BaseListerProps<T>) {
  const [selected, setSelected] = useState('');
  const {state, items, retry} = useResourceStream<T & WithMeta>(objectType);

  const handleSelect = (item: T | null) => {
    if (item === null || Object.keys(item as object).length === 0) {
      setSelected('');
      return;
    }
    setSelected((item as WithMeta)?.metadata?.name || '');
  };

  const filteredItem: T | undefined = selected
    ? (FilterItems(selected, items) as T | undefined)
    : undefined;

  const lower = title.toLowerCase();

  const body = () => {
    if (state === 'connecting') {
      return <GridCol span={12}><CenteredLoader/></GridCol>;
    }
    if (state === 'error') {
      return (
        <GridCol span={12}>
          <ErrorState message={`Unable to load ${lower}. The connection may be unavailable.`} onRetry={retry}/>
        </GridCol>
      );
    }
    if (filteredItem) {
      return renderDetails(filteredItem);
    }
    if (state === 'empty') {
      return <GridCol span={12}><EmptyState label={`No ${lower} found`}/></GridCol>;
    }
    return renderTable(items as T[], handleSelect);
  };

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol span={{base: 12, sm: 9}}>
        <Title className={sourceCodePro400.className} order={2}>
          {title}
        </Title>
      </GridCol>
      <GridCol span={{base: 12, sm: 3}} className="flex justify-end items-center">
        {selected &&
          <IconArrowBigLeft
            onClick={() => handleSelect(null)}
            size={32}
            className="cursor-pointer hover:opacity-70"
            role="button"
            tabIndex={0}
            aria-label="Back to list"
          />
        }
      </GridCol>
      {body()}
    </Grid>
  );
}
