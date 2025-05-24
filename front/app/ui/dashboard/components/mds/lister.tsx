'use client';

import Link from 'next/link';
import React, { useState, useEffect } from 'react';
import {sourceCodePro400} from "@/fonts";

import Search from '@/app/ui/dashboard/search'
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title } from '@mantine/core';

import MachineDeploymentTable from '@/app/ui/dashboard/components/mds/table'
import MachineDeploymentDetails from "@/app/ui/dashboard/components/mds/details";

import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";
import BaseLister from "@/app/ui/dashboard/base/lister";
import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import MachineDetails from "@/app/ui/dashboard/components/machines/details";
import MachinesTable from "@/app/ui/dashboard/components/machines/table";

/**
 * MachineDeploymentLister is a functional component responsible for listing machine deployments
 * and optionally showing details for a selected deployment. It manages state for the deployments,
 * filters based on user selection, and handles WebSocket communication for data fetching.
 */
export default function MachineDeploymentLister() {
  return <BaseLister
    objectType="machine-deployment"
    items={[]}
    renderDetails={(item: MachineDeploymentType) => <MachineDeploymentDetails md={item}/>}
    renderTable={(items : MachineDeploymentType[], handleSelect) =>  (
      <MachineDeploymentTable  mds={items}/>
    )}
    title="Machine Deployments / cluster.x-k8s.io"
  />
}
