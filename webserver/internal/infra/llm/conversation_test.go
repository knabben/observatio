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
				stopper:  tt.keepCount,
				messages: make([]anthropic.MessageParam, tt.initialLength),
			}
			cm.TrimHistory()
			assert.Equal(t, tt.expectedCount, len(cm.messages))
		})
	}
}

func TestGetConversationHistory(t *testing.T) {
	tests := []struct {
		name          string
		initialLength int
		stopper       int
		expectLength  int
	}{
		{
			name:          "EmptyHistory",
			initialLength: 0,
			stopper:       5,
			expectLength:  0,
		},
		{
			name:          "UnderStopper",
			initialLength: 3,
			stopper:       5,
			expectLength:  3,
		},
		{
			name:          "EqualToStopper",
			initialLength: 5,
			stopper:       5,
			expectLength:  5,
		},
		{
			name:          "OverStopper",
			initialLength: 10,
			stopper:       7,
			expectLength:  7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConversationManager{
				stopper:  tt.stopper,
				messages: make([]anthropic.MessageParam, tt.initialLength),
			}
			history := cm.GetConversationHistory()
			assert.Equal(t, tt.expectLength, len(history))
		})
	}
}
