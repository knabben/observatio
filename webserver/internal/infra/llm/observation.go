package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"

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
	ID        string                 `json:"id"`
	AgentID   string                 `json:"agent_id"`
	Content   string                 `json:"content"`
	Type      string                 `json:"type"` // "user", "agent", "system"
	Timestamp time.Time              `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
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

	// Process response and handle tool calls
	responseText := ""
	var toolResults []string

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

	chatMessage := &ChatMessage{
		ID:        generateID(),
		Content:   responseText,
		Type:      "agent",
		Timestamp: time.Now(),
		Context: map[string]interface{}{
			"agent_type": "agent",
		},
	}

	return chatMessage, nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
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
		{
			ID:           "remediation-agent",
			Name:         "Remediation Agent",
			Type:         "execution",
			Status:       "standby",
			Activity:     "Ready to execute automated fixes",
			Capabilities: []string{"auto-remediation", "workflow-execution", "safety-checks"},
			LastUpdate:   time.Now(),
		},
		{
			ID:           "monitoring-agent",
			Name:         "Monitoring Agent",
			Type:         "monitoring",
			Status:       "active",
			Activity:     "Continuous health monitoring",
			Capabilities: []string{"metrics-collection", "anomaly-detection", "alerting"},
			LastUpdate:   time.Now(),
		},
	}

	for _, agent := range agents {
		s.agents[agent.ID] = agent
	}
}
