package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
	// anthropicClient is the client used for interacting with the Anthropic API.
	anthropicClient Client

	// agents is a map that stores AI agents identified by their unique string IDs.
	agents map[string]*Agent

	// wsConnections maps agent IDs to active WebSocket connections for real-time communication.
	wsConnections map[string]*websocket.Conn

	// chatHistory maps agent IDs to their respective chat message history for maintaining conversational context.
	chatHistory map[string][]ChatMessage

	// tools represent a collection of tools available for the ObservationService to execute specific operations or commands.
	tools []Tool
}

func NewObservationService(client Client) (*ObservationService, error) {
	service := &ObservationService{
		anthropicClient: client,
		agents:          make(map[string]*Agent),
		wsConnections:   make(map[string]*websocket.Conn),
		chatHistory:     make(map[string][]ChatMessage),
		tools:           initializeTools(),
	}
	service.initializeAgents()

	return service, nil
}

func (s *ObservationService) ChatWithAgent(ctx context.Context, message ChatMessage, agentID string) (*ChatMessage, error) {
	logger := log.FromContext(ctx)
	//client := s.anthropicClient.GetClient()

	messages := []anthropic.MessageParam{
		{
			Content: []anthropic.ContentBlockParamUnion{
				{OfRequestTextBlock: &anthropic.TextBlockParam{Text: formatMessage(message.Content)}},
			},
			Role: anthropic.MessageParamRoleAssistant,
		},
	}

	// Add chat conversation history for the context
	if history, exists := s.chatHistory[agentID]; !exists {
		s.chatHistory[agentID] = []ChatMessage{}
	} else if len(history) > 0 {
		lastElements := getLastElements(history, 5)
		logger.Info("Found history messages", "messages", len(lastElements))
		for i := len(lastElements) - 1; i >= 0; i-- {
			role := anthropic.MessageParamRoleUser
			text := formatMessage(history[i].Content)
			if history[i].Actor == "agent" {
				role = anthropic.MessageParamRoleAssistant
				text = history[i].Content
			}
			messages = append(messages, anthropic.MessageParam{
				Content: []anthropic.ContentBlockParamUnion{
					{OfRequestTextBlock: &anthropic.TextBlockParam{Text: text}},
				},
				Role: role,
			})
		}
	}

	logger.Info("messages history", "history", messages)
	request := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 4000,
		System: []anthropic.TextBlockParam{
			{Text: TASK_SYSTEM},
		},
		Messages: messages,
	}

	logger.Info("Request to Claude", "request", len(request.Messages))
	for _, res := range messages {
		logger.Info("Request to Claude", "request", res)
	}
	//response, err := client.Messages.New(ctx, request)
	//if err != nil {
	//	return nil, fmt.Errorf("claude API error: %v", err)
	//}

	var toolResults []string
	var responseText string
	//logger.Info("Response from Claude", "response", response)
	//for _, content := range response.Content {
	//	switch content := content.AsAny().(type) {
	//	case anthropic.TextBlock:
	//		responseText += content.Text
	//	}
	//}

	// Combine response text with tool results
	if len(toolResults) > 0 {
		responseText += "\n\nTool Results:\n" + fmt.Sprintf("%v", toolResults)
	}

	botMessage := ChatMessage{
		ID:        generateID(),
		Content:   strings.ReplaceAll("response random", "\n", "<br />"),
		Type:      "chatbot",
		Actor:     "agent",
		AgentID:   "cluster-agent",
		Timestamp: time.Now().Format("01/02/2006 15:04:05"),
	}
	s.chatHistory[agentID] = append(s.chatHistory[agentID], []ChatMessage{message, botMessage}...)
	return &botMessage, nil
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

func getLastElements(s []ChatMessage, n int) []ChatMessage {
	if len(s) < n {
		return s
	}
	return s[len(s)-n:]
}

func reverseElements(result []ChatMessage) []ChatMessage {
	for i := 0; i < len(result)/2; i++ {
		j := len(result) - 1 - i
		result[i], result[j] = result[j], result[i]
	}
	return result
}

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
