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

When you receive a customer question, you must respond with a detailed DESCRIPTION of the issue, under <description></description> tag, ALWAYS.
If you have suggestions for fixing the issue or improvements, you can also include them under <suggestions></suggestions> tag.

You have native Kubernetes Tools to execute commands as part of your troubleshooting capabilities:
 - kubectl commands: Execute cluster inspection and modification commands
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
}

func ToMessageParam(message string) *ChatMessage {
	var (
		messageType      = "chatbot"
		messageActor     = "agent"
		agentID          = "cloud-agent"
		messageTimestamp = time.Now().Format("01/02/2006 15:04:05")
	)

	return &ChatMessage{
		ID:        generateID(),
		Content:   strings.ReplaceAll(message, "\n", "<br />"),
		Type:      messageType,
		Actor:     messageActor,
		AgentID:   agentID,
		Timestamp: messageTimestamp,
	}
}

// formatMessage creates a formatted message combining the task context, error details,
// and expected output format for the LLM processing
func formatMessage(errorMessage string) string {
	var messageBuilder strings.Builder
	formattedQuestion := fmt.Sprintf(questionFormat, errorMessage)
	messageBuilder.WriteString(fmt.Sprintf(messageTemplate, TASK_CONTEXT, formattedQuestion))
	return messageBuilder.String()
}
