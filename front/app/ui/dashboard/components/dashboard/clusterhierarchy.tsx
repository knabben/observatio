'use client'


import React, {useState,useEffect,useCallback} from "react";
import {roboto} from "@/fonts";
import Header from "@/app/ui/dashboard/utils/header";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {getClusterHierarchy} from "@/app/lib/data";
import {Card, Text} from "@mantine/core";

import {
  ReactFlow,
    addEdge,
    useEdgesState,
    useNodesState,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

type Hierarchy = {
  nodes: any[],
  edges: any[],
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
//
// export const initialNodes = [
//   {
//     id: '1',
//     data: { label: 'Node 1' },
//     position: { x: 150, y: 0 },
//     style: { backgroundColor: '#6ede87', color: '#000000' },
//   },
//   {
//     id: '2',
//     data: { label: 'Node 2' },
//     position: { x: 0, y: 150 },
//   },
//   {
//     id: '3',
//     data: { label: 'Node 3' },
//     position: { x: 300, y: 150 },
//   },
//
// ];
// export const initialEdges = [
//   { id: 'e1-2', source: '1', target: '2' },
//   { id: 'e1-3', source: '4', target: '2' },
//   { id: 'e1-4', source: '3', target: '1' },
// ];
export function RenderTopology({
  hierarchy,
}: {hierarchy: Hierarchy}) {
  const [nodes, setNodes, onNodesChange] = useNodesState(hierarchy?.nodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(hierarchy?.edges);

  const onConnect = useCallback(
    (connection: any) => setEdges((eds) => addEdge(connection, eds)),
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
            // @ts-ignore
            !isLoading  && <RenderTopology hierarchy={hierarchy} />
          }
      </div>
    </Card>
  );
}
