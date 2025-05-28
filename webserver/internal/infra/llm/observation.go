package llm

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	anthropicClient anthropic.Client
	k8sClient       kubernetes.Interface
	agents          map[string]*Agent
	chatHistory     map[string][]ChatMessage
	wsConnections   map[string]*websocket.Conn
	mu              sync.RWMutex
	tools           []Tool
}

func NewObservationService() (*ObservationService, error) {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config: %v", err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %v", err)
	}

	service := &ObservationService{
		anthropicClient: NewClient(),
		k8sClient:       k8sClient,
		agents:          make(map[string]*Agent),
		wsConnections:   make(map[string]*websocket.Conn),
		tools:           initializeTools(),
	}

	// Initialize agents
	service.initializeAgents()

	return service, nil
}

func (s *ObservationService) buildSystemPrompt(agent *Agent) string {
	return fmt.Sprintf(`You are %s, a specialized AI agent in the Observatio platform for Kubernetes cluster management and troubleshooting.

Your role: %s
Your capabilities: %v

You have access to various tools for Kubernetes cluster analysis and remediation. Use these tools when appropriate to provide accurate, actionable insights.

Always be specific, technical, and provide actionable recommendations. Focus on:
1. Clear problem identification
2. Root cause analysis
3. Step-by-step remediation plans
4. Preventive measures

Context awareness: You can access cluster state, metrics, logs, and historical data through the available tools.`,
		agent.Name, agent.Type, agent.Capabilities)
}

func (s *ObservationService) ChatWithAgent(ctx context.Context, agentID, message string) (*ChatMessage, error) {
	agent, exists := s.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not found", agentID)
	}

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(message)),
	}

	request := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 4000,
		System:    []anthropic.TextBlockParam{
			{Text: s.buildSystemPrompt(agent)},
		}
		Messages:  messages,
	}

	// Call Claude API
	response, err := s.anthropicClient.Messages.New(ctx, request)
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
		AgentID:   agentID,
		Content:   responseText,
		Type:      "agent",
		Timestamp: time.Now(),
		Context: map[string]interface{}{
			"agent_type": agent.Type,
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
