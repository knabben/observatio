import React from "react";
import {Group, Indicator, Text} from "@mantine/core";
import {StatusState} from "@/app/ui/dashboard/shared/status";

interface StatusIndicatorProps {
  state: StatusState;
  /** Accessible + visible label; defaults per state. */
  label?: string;
  size?: number;
  /** Render only the dot (no text label). */
  dotOnly?: boolean;
}

const CONFIG: Record<StatusState, {color: string; label: string}> = {
  healthy: {color: "green", label: "Ready"},
  notready: {color: "red", label: "Not ready"},
  unknown: {color: "gray", label: "Unknown"},
};

/**
 * Tri-state status indicator. Distinguishes healthy / not-ready / unknown by color AND
 * an accessible label (never color alone), and NEVER animates — a static/failed state
 * must not imply work-in-progress.
 */
export const StatusIndicator: React.FC<StatusIndicatorProps> = ({
  state,
  label,
  size = 14,
  dotOnly = false,
}) => {
  const {color, label: defaultLabel} = CONFIG[state];
  const text = label ?? defaultLabel;
  const dot = (
    <Indicator inline color={color} size={size} aria-label={text} role="img" processing={false}/>
  );
  if (dotOnly) return dot;
  return (
    <Group gap="xs" wrap="nowrap">
      {dot}
      <Text size="sm">{text}</Text>
    </Group>
  );
};
