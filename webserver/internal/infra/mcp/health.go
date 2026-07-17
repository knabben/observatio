package mcp

import (
	"context"
	"time"
)

// healthCheckInterval is how often an external source's health is re-probed once it's connected
// (research.md R4) - out-of-band, never inline with a chat turn.
const healthCheckInterval = 30 * time.Second

// StartHealthChecking launches a background goroutine that re-probes this source on a fixed
// interval, refreshing its cached HealthStatus and capability list, until ctx is done. It returns
// immediately. NewMCPToolSource already performs the first probe synchronously, so this loop's
// first tick is a re-check, not the initial one.
func (s *MCPToolSource) StartHealthChecking(ctx context.Context) {
	go runHealthLoop(ctx, healthCheckInterval, s.probe, s.recordProbeResult)
}

// runHealthLoop calls probe on the given interval until ctx is done, passing each result to
// record. Extracted from StartHealthChecking so the health-check *timing* behavior (state
// transitions across ticks) can be tested with a short interval instead of the real 30s one.
func runHealthLoop(ctx context.Context, interval time.Duration, probe func(context.Context) error, record func(error)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			record(probe(ctx))
		}
	}
}
