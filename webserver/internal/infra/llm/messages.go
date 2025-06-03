package llm

import (
	"fmt"
	"strings"
)

var (
	TASK_SYSTEM  = `You will serve as a Kubernetes administrator managing a on-premises datacenter on VMware vCenter.`
	TASK_CONTEXT = `Your task is to assist operators in troubleshooting issues within the cluster.
Provide a detailed explanation of the issue. New inputs are provided and you must respond accordingly the context
Do not use any existent tool for now.`
)

// MessageTemplate defines the structure for formatting error messages
const (
	messageTemplate = "%s\n%s"
	questionFormat  = "Here is the customer question: <question>%s</question>"
)

type ChatMessage struct {
	// ID represents the unique identifier for the chat message.
	ID string `json:"id"`

	// AgentID represents the identifier of the agent associated with the chat message.
	AgentID string `json:"agent_id"`

	// Content represents the text content of the chat message.
	Content string `json:"content"`

	// Type represents the category or role of the chat message, typically denoting its origin or purpose.
	Type string `json:"type"`

	// Actor specifies the origin of the message, such as "agent" or "user".
	Actor string `json:"actor"`

	// Timestamp represents the time when the chat message was created or sent.
	Timestamp string `json:"timestamp"`
}

// formatMessage creates a formatted message combining the task context, error details,
// and expected output format for the LLM processing
func formatMessage(errorMessage string) string {
	var messageBuilder strings.Builder
	formattedQuestion := fmt.Sprintf(questionFormat, errorMessage)
	messageBuilder.WriteString(fmt.Sprintf(messageTemplate, TASK_CONTEXT, formattedQuestion))
	return messageBuilder.String()
}
