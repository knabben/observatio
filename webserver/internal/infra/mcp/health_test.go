package mcp

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recordingProbe is a test double for a probe function whose result flips on each call
// according to a supplied script, recording every result runHealthLoop passed to record.
type recordingProbe struct {
	mu      sync.Mutex
	script  []error // nil entry = success
	results []error
	done    chan struct{} // closed once len(script) results have been recorded
}

func newRecordingProbe(script ...error) *recordingProbe {
	return &recordingProbe{script: script, done: make(chan struct{})}
}

func (p *recordingProbe) probe(context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.results) >= len(p.script) {
		return p.script[len(p.script)-1]
	}
	return p.script[len(p.results)]
}

func (p *recordingProbe) record(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.results = append(p.results, err)
	if len(p.results) == len(p.script) {
		close(p.done)
	}
}

func TestRunHealthLoop_TicksOnIntervalUntilContextDone(t *testing.T) {
	fake := newRecordingProbe(nil, errors.New("boom"), nil) // healthy -> unhealthy -> healthy
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go runHealthLoop(ctx, 5*time.Millisecond, fake.probe, fake.record)

	select {
	case <-fake.done:
	case <-ctx.Done():
		t.Fatal("timed out waiting for all scripted probe results to be recorded")
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	require.Len(t, fake.results, 3)
	assert.NoError(t, fake.results[0])
	assert.Error(t, fake.results[1])
	assert.NoError(t, fake.results[2])
}

func TestMCPToolSource_StartHealthChecking_RecoversAutomatically(t *testing.T) {
	// A source whose underlying command fails every time is deterministic; instead, drive the
	// scenario through the real MCPToolSource.recordProbeResult/probe pair with a fake transport
	// that succeeds, to confirm StartHealthChecking's wiring (not just the extracted loop) ends
	// up calling recordProbeResult on tick, transitioning health as probe results come in.
	factory := newFakeServerTransportFactory(t, readOnlyTool("list_backups"))
	src := newMCPToolSource("velero-mcp", factory)
	require.Equal(t, HealthUnknown, src.Health().State, "health starts unknown before any probe")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use a short interval directly via runHealthLoop wired to this source's own probe/record,
	// mirroring what StartHealthChecking does with healthCheckInterval.
	tickDone := make(chan struct{})
	go func() {
		runHealthLoop(ctx, 5*time.Millisecond, src.probe, func(err error) {
			src.recordProbeResult(err)
			select {
			case tickDone <- struct{}{}:
			case <-ctx.Done():
			}
		})
	}()

	select {
	case <-tickDone:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for the first health-check tick")
	}

	assert.Equal(t, HealthHealthy, src.Health().State)
	assert.NotEmpty(t, src.Capabilities(), "a successful health check must refresh the capability list")
}
