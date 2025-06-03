package llm

import (
	"github.com/anthropics/anthropic-sdk-go"
)

// ConversationManager handles message history
type ConversationManager struct {
	client   *anthropic.Client
	messages []anthropic.MessageParam
}

// NewConversationManager creates a new conversation manager
func NewConversationManager(client *anthropic.Client) *ConversationManager {
	return &ConversationManager{
		client:   client,
		messages: make([]anthropic.MessageParam, 0),
	}
}

// AddUserMessage adds a user message to the conversation history
func (cm *ConversationManager) AddUserMessage(content string) {
	message := anthropic.NewUserMessage(anthropic.NewTextBlock(content))
	cm.messages = append(cm.messages, message)
}

// AddAssistantMessage adds an assistant message to the conversation history
func (cm *ConversationManager) AddAssistantMessage(content string) {
	message := anthropic.NewAssistantMessage(anthropic.NewTextBlock(content))
	cm.messages = append(cm.messages, message)
}

// GetConversationHistory returns the current conversation history
func (cm *ConversationManager) GetConversationHistory() []anthropic.MessageParam {
	return cm.messages
}

// ClearHistory clears the conversation history
func (cm *ConversationManager) ClearHistory() {
	cm.messages = make([]anthropic.MessageParam, 0)
}

// GetHistoryLength returns the number of messages in history
func (cm *ConversationManager) GetHistoryLength() int {
	return len(cm.messages)
}

// TrimHistory removes older messages to stay within token limits
// Keeps the last 'keepCount' messages
func (cm *ConversationManager) TrimHistory(keepCount int) {
	if len(cm.messages) < keepCount {
		return
	}
	if len(cm.messages) > keepCount {
		cm.messages = cm.messages[len(cm.messages)-keepCount:]
	}
}
