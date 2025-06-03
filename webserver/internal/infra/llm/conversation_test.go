package llm

import (
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestTrimHistory(t *testing.T) {
	tests := []struct {
		name          string
		initialLength int
		keepCount     int
		expectedCount int
	}{
		{
			name:          "NoTrimmingNeeded",
			initialLength: 3,
			keepCount:     5,
			expectedCount: 3,
		},
		{
			name:          "ExactKeepCount",
			initialLength: 5,
			keepCount:     5,
			expectedCount: 5,
		},
		{
			name:          "TrimToKeepCount",
			initialLength: 10,
			keepCount:     7,
			expectedCount: 7,
		},
		{
			name:          "KeepNone",
			initialLength: 5,
			keepCount:     0,
			expectedCount: 0,
		},
		{
			name:          "KeepOne",
			initialLength: 5,
			keepCount:     1,
			expectedCount: 1,
		},
		{
			name:          "EmptyHistory",
			initialLength: 0,
			keepCount:     5,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConversationManager{
				client:   &anthropic.Client{},
				messages: make([]anthropic.MessageParam, tt.initialLength),
			}
			cm.TrimHistory(tt.keepCount)
			assert.Equal(t, tt.expectedCount, len(cm.messages))
		})
	}
}
