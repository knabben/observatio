'use client'

import React, {useState,useEffect,useCallback} from "react";
import Header from "@/app/ui/dashboard/utils/header";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {getClusterHierarchy} from "@/app/lib/data";
import {Card, Text} from "@mantine/core";

import {
  ReactFlow,
  addEdge,
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
  const [hierarchy, setHierarchy] = useState<Hierarchy>();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleFetchError = (error: Error) => {
    console.error('Failed to fetch cluster summary:', error);
    setError('Failed to load cluster summary');
    setIsLoading(false);
  };

  useEffect(() => {
    const fetchClusterTopology= async () => {
      try {
        setIsLoading(true);
        const response = await getClusterHierarchy();
        setHierarchy(response);
      } catch (error) {
        handleFetchError(error as Error)
      } finally {
        setIsLoading(false);
      }
    };
    fetchClusterTopology();
  }, []);
  return {hierarchy, isLoading, error};
};

export function RenderTopology({
  hierarchy,
}: {hierarchy: Hierarchy}) {
  const [nodes, , onNodesChange] = useNodesState(hierarchy?.nodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(hierarchy?.edges);

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
    />
  )
}
export default function ClusterHierarchy() {
  const {hierarchy, isLoading, error} = useClusterHierarchy();
  return (
    <Card shadow="md" radius="md" withBorder className="text-center" >
      <div style={{ width: '860px', height: '500px' }}>
        <Header title="Cluster Topology"/>
          {isLoading && <CenteredLoader />}
          {error && <Text c="red">{error}</Text>}
          {
            // @ts-expect-error undefined is not assignable to type 'false'
            !isLoading  && <RenderTopology hierarchy={hierarchy} />
          }
      </div>
    </Card>
  );
}
