package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Agent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Activity     string    `json:"activity"`
	LastUpdate   time.Time `json:"last_update"`
	Capabilities []string  `json:"capabilities"`
}

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

type ObservationService struct {
	anthropicClient Client
	agents          map[string]*Agent
	wsConnections   map[string]*websocket.Conn
	tools           []Tool
}

func NewObservationService(client Client) (*ObservationService, error) {
	service := &ObservationService{
		anthropicClient: client,
		agents:          make(map[string]*Agent),
		wsConnections:   make(map[string]*websocket.Conn),
		tools:           initializeTools(),
	}
	service.initializeAgents()

	return service, nil
}

func (s *ObservationService) ChatWithAgent(ctx context.Context, message string) (*ChatMessage, error) {
	client := s.anthropicClient.GetClient()
	request := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 4000,
		System: []anthropic.TextBlockParam{
			{Text: TASK_SYSTEM},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(formatMessage(message))),
		},
	}

	response, err := client.Messages.New(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("claude API error: %v", err)
	}

	//Process response and handle tool calls
	var toolResults []string
	var responseText string
	for _, content := range response.Content {
		switch content := content.AsAny().(type) {
		case anthropic.TextBlock:
			responseText += content.Text
		}
	}

	// Combine response text with tool results
	if len(toolResults) > 0 {
		responseText += "\n\nTool Results:\n" + fmt.Sprintf("%v", toolResults)
	}

	return &ChatMessage{
		ID:        generateID(),
		Content:   strings.ReplaceAll(responseText, "\n", "<br />"),
		Type:      "chatbot",
		Actor:     "agent",
		AgentID:   "cluster-agent",
		Timestamp: time.Now().Format("01/02/2006 15:04:05"),
	}, nil
}

func generateID() string {
	return uuid.New().String()
}

func initializeTools() []Tool {
	return []Tool{
		{
			Name:        "kubectl_command",
			Description: "Execute kubectl commands to inspect or modify Kubernetes resources",
			InputSchema: ToolSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The kubectl command to execute",
					},
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (optional)",
					},
					"dry_run": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to run in dry-run mode",
					},
				},
				Required: []string{"command"},
			},
		},
	}
}

// Initialize AI agents
func (s *ObservationService) initializeAgents() {
	agents := []*Agent{
		{
			ID:           "cluster-agent",
			Name:         "Cluster Agent",
			Type:         "analysis",
			Status:       "active",
			Activity:     "Monitoring cluster health and analyzing issues",
			Capabilities: []string{"cluster-analysis", "resource-monitoring", "pattern-detection"},
			LastUpdate:   time.Now(),
		},
	}

	for _, agent := range agents {
		s.agents[agent.ID] = agent
	}
}
