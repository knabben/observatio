import React from "react";
import {Group, Indicator, Text} from "@mantine/core";
import {StatusState} from "@/app/ui/dashboard/shared/status";
import {STATUS_COLORS} from "@/app/styles/theme";

interface StatusIndicatorProps {
  state: StatusState;
  /** Accessible + visible label; defaults per state. */
  label?: string;
  size?: number;
  /** Render only the dot (no text label). */
  dotOnly?: boolean;
}

const CONFIG: Record<StatusState, {color: string; label: string}> = {
  healthy: {color: STATUS_COLORS.healthy, label: "Ready"},
  degraded: {color: STATUS_COLORS.degraded, label: "Degraded"},
  notready: {color: STATUS_COLORS.notready, label: "Not ready"},
  unknown: {color: STATUS_COLORS.unknown, label: "Unknown"},
};

/**
 * Status indicator. Distinguishes healthy / degraded / not-ready / unknown by color AND
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
