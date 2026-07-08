package llm

import (
	"fmt"
	"strings"
	"time"
)

var (
	TASK_SYSTEM  = `You will serve as a Cluster API advisor helping troubleshoot on-premises Kubernetes on VMware vCenter infrastructure.`
	TASK_CONTEXT = `You are Observatio AI, an advanced Kubernetes cluster troubleshooting assistant specialized in ClusterAPI (CAPI) environments. 
You are part of a sophisticated monitoring platform that revolutionizes Kubernetes cluster assessment through data aggregation, MCP server integration, and AI-powered remediation.
Provide expert-level Kubernetes troubleshooting, root cause analysis, and automated remediation for DevOps teams managing distributed CAPI clusters. 
You excel at translating complex cluster issues into actionable solutions. Your task is to assist operators in troubleshooting issues within the cluster.

When you receive a customer question, always respond with a detailed description of the issue in plain prose.
If you have suggestions for fixing the issue or improvements, include them under a "Suggestions:" heading.
Do not use XML/HTML-style tags in your response (e.g. no <description> or <suggestions> tags) - the response is
rendered as plain text, so any such tags would show up literally instead of being formatted.

You only run available tools when requested to increase the context of the issue.
`
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

	// Event distinguishes a one-shot complete message (empty, the default) from a streamed
	// "delta" chunk to append to the message with the same ID, or a "done" chunk marking the
	// end of a streamed reply.
	Event string `json:"event,omitempty"`
}

func timestamp() string {
	return time.Now().Format("01/02/2006 15:04:05")
}

func ToMessageParam(message string) *ChatMessage {
	return &ChatMessage{
		ID:        generateID(),
		Content:   message,
		Type:      "chatbot",
		Actor:     "agent",
		AgentID:   "cloud-agent",
		Timestamp: timestamp(),
	}
}

// newStreamChunk builds a streamed chatbot chunk sharing responseID across a whole reply, so the
// frontend can append "delta" chunks to the right in-progress message and knows when it's "done".
func newStreamChunk(responseID, content, event string) *ChatMessage {
	return &ChatMessage{
		ID:        responseID,
		Content:   content,
		Type:      "chatbot",
		Actor:     "agent",
		AgentID:   "cloud-agent",
		Timestamp: timestamp(),
		Event:     event,
	}
}

// formatMessage creates a formatted message combining the task context, error details,
// and expected output format for the LLM processing
func formatMessage(errorMessage string) string {
	var messageBuilder strings.Builder
	formattedQuestion := fmt.Sprintf(questionFormat, errorMessage)
	fmt.Fprintf(&messageBuilder, messageTemplate, TASK_CONTEXT, formattedQuestion)
	return messageBuilder.String()
}
