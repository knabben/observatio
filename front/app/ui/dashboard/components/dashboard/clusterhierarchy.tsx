'use client'

import React, {useCallback} from "react";
import Header from "@/app/ui/dashboard/utils/header";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {getClusterHierarchy} from "@/app/lib/data";
import {Card, Text} from "@mantine/core";
import {EmptyState} from "@/app/ui/dashboard/shared/empty-state";
import {useFetchState} from "@/app/ui/dashboard/shared/use-fetch-state";

import {
  ReactFlow,
  addEdge,
  Background,
  Controls,
  useEdgesState,
  useNodesState, Connection,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

type Node = {
  id: string,
  type: string,
  data: {
    label: string,
  },
  position: {
    x: number,
    y: number,
  },
  style: {
    background: string,
    color: string,
    border: string,
  }
}

type Edge = {
  id: string,
  source: string,
  target: string,
  type: string,
  animated: boolean,
}

type Hierarchy = {
  nodes: Node[],
  edges: Edge[],
}

export const useClusterHierarchy = () => {
  const {data: hierarchy, isLoading, error} = useFetchState<Hierarchy | undefined>(
    getClusterHierarchy,
    undefined,
    'Failed to load cluster topology',
  );
  return {hierarchy, isLoading, error};
};

export function RenderTopology({
  hierarchy,
}: {hierarchy: Hierarchy}) {
  const [nodes, , onNodesChange] = useNodesState(hierarchy.nodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(hierarchy.edges);

  const onConnect = useCallback(
    (connection: Connection) => setEdges((eds) => addEdge(connection, eds)),
    [setEdges],);

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      fitView
    >
      <Background/>
      <Controls/>
    </ReactFlow>
  )
}
export default function ClusterHierarchy() {
  const {hierarchy, isLoading, error} = useClusterHierarchy();
  const isEmpty = !isLoading && !error && (hierarchy?.nodes?.length ?? 0) === 0;
  return (
    <Card shadow="md" radius="md" withBorder className="text-center" >
      <div style={{ width: '100%', height: '500px' }}>
        <Header title="Cluster Topology"/>
          {isLoading && <CenteredLoader />}
          {error && <Text c="red">{error}</Text>}
          {isEmpty && <EmptyState label="No cluster topology found"/>}
          {!isLoading && !error && !isEmpty && hierarchy &&
            <RenderTopology hierarchy={hierarchy} />
          }
      </div>
    </Card>
  );
}
