'use client'


import React, {useState,useEffect,useCallback} from "react";
import {roboto} from "@/fonts";
import Header from "@/app/ui/dashboard/utils/header";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {getClusterClasses} from "@/app/lib/data";
import {Card} from "@mantine/core";

import {
  ReactFlow,
    addEdge,
    useEdgesState,
    useNodesState,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

export const useClusterHierarchy = () => {
  const [hierarchy, setHierarchy] = useState<[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleFetchError = (error: Error) => {
    console.error('Failed to fetch cluster summary:', error);
    setError('Failed to load cluster summary');
    setIsLoading(false);
  };

  useEffect(() => {
    const fetchClusterClasses = async () => {
      try {
        setIsLoading(true);
        // const response = await getClusterClasses();
        // setClusterClasses(response);
      } catch (error) {
        handleFetchError(error as Error)
      } finally {
        setIsLoading(false);
      }
    };
    fetchClusterClasses();
  }, []);

  return {hierarchy, isLoading, error};
};

export const initialNodes = [
  {
    id: '1',
    data: { label: 'Node 1' },
    position: { x: 150, y: 0 },
    style: { backgroundColor: '#6ede87', color: '#000000' },
  },
  {
    id: '2',
    data: { label: 'Node 2' },
    position: { x: 0, y: 150 },
  },
  {
    id: '3',
    data: { label: 'Node 3' },
    position: { x: 300, y: 150 },
  },
];
export const initialEdges = [
  { id: 'e1-2', source: '1', target: '2' },
  { id: 'e1-3', source: '1', target: '3',  animated: true  },
];
export default function ClusterHierarchy() {
  const {hierarchy, isLoading, error} = useClusterHierarchy();
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const onConnect = useCallback(
    (connection) => setEdges((eds) => addEdge(connection, eds)),
    [setEdges],
  );

  return (
    <Card shadow="md" radius="md" withBorder className="text-center" >
      <div style={{ width: '800px', height: '500px' }}>
        <Header title="Cluster Topology"/>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
        />
      </div>
    </Card>
  );
}
