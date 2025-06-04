package llm

import (
	"github.com/anthropics/anthropic-sdk-go"
)

// ConversationManager handles message history
type ConversationManager struct {
	stopper  int
	messages []anthropic.MessageParam
}

// NewConversationManager creates a new conversation manager
func NewConversationManager(stopper int) *ConversationManager {
	return &ConversationManager{
		stopper:  stopper,
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
	if len(cm.messages) <= cm.stopper {
		return cm.messages
	}
	return cm.messages[:cm.stopper]
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
// Keeps the last 'stopper' messages
func (cm *ConversationManager) TrimHistory() {
	if len(cm.messages) <= cm.stopper {
		return
	}
	cm.messages = cm.messages[len(cm.messages)-cm.stopper:]
}
