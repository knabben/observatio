package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
	tools []anthropic.ToolUnionParam
}

func NewObservationService(client Client) (*ObservationService, error) {
	allTools := RenderTools()
	tools := make([]anthropic.ToolUnionParam, len(allTools))
	for i, toolParam := range allTools {
		tools[i] = anthropic.ToolUnionParam{OfTool: &toolParam}
	}

	service := &ObservationService{
		anthropicClient: client,
		agents:          make(map[string]*Agent),
		wsConnections:   make(map[string]*websocket.Conn),
		chatHistory:     make(map[string][]ChatMessage),
		tools:           tools,
	}
	service.initializeAgents()

	return service, nil
}

func (s *ObservationService) ChatWithAgent(ctx context.Context, message ChatMessage, agentID string) (*ChatMessage, error) {
	logger := log.FromContext(ctx)
	client := s.anthropicClient.GetClient()

	messages := []anthropic.MessageParam{
		{
			Content: []anthropic.ContentBlockParamUnion{
				{OfRequestTextBlock: &anthropic.TextBlockParam{Text: formatMessage(message.Content)}},
			},
			Role: anthropic.MessageParamRoleUser,
		},
	}

	// Add chat conversation history for the context
	if history, exists := s.chatHistory[agentID]; !exists {
		s.chatHistory[agentID] = []ChatMessage{}
	} else if len(history) > 0 {
		lastElements := getLastElements(history, 2)
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

	request := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 4000,
		System: []anthropic.TextBlockParam{
			{Text: TASK_SYSTEM},
		},
		Messages: messages,
		Tools:    s.tools,
	}

	response, err := client.Messages.New(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("claude API error: %v", err)
	}

	var responseText string
	var toolResults []string

	logger.Info("Response from Claude", "response", response)
	for _, block := range response.Content {
		switch content := block.AsAny().(type) {
		case anthropic.TextBlock:
			responseText += content.Text

		case anthropic.ToolUseBlock:
			logger.Info("Tools use command")
			// save as part of message
			inputJSON, _ := json.Marshal(block.Input)
			logger.Info(block.Name + ": " + string(inputJSON))

			var response interface{}
			switch block.Name {
			case "kubectl":
				var input struct {
					Command string `json:"command"`
				}

				err := json.Unmarshal([]byte(block.JSON.Input.Raw()), &input)
				if err != nil {
					panic(err)
				}

				response, err = RunKubectl(input.Command)
				if err != nil {
					logger.Error(err, "Error running kubectl command")
					response = err.Error()
				}
				logger.Info("Kubectl response", "response", response)
				toolResults = append(toolResults, response.(string))
			}
		}
	}

	// Combine response text with tool results
	if len(toolResults) > 0 {
		responseText += "\n\nTool Results:\n" + fmt.Sprintf("%v", toolResults)
	}

	botMessage := ChatMessage{
		ID:        generateID(),
		Content:   strings.ReplaceAll(responseText, "\n", "<br />"),
		Type:      "chatbot",
		Actor:     "agent",
		AgentID:   "cluster-agent",
		Timestamp: time.Now().Format("01/02/2006 15:04:05"),
	}
	s.chatHistory[agentID] = append(s.chatHistory[agentID], []ChatMessage{message, botMessage}...)
	return &botMessage, nil
}

func RunKubectl(command string) (string, error) {
	// Execute kubectl command using os/exec
	cmd := exec.Command("kubectl", strings.Fields(command)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing kubectl command: %v", err)
	}

	return string(output), nil
}

func generateID() string {
	return uuid.New().String()
}

func getLastElements(s []ChatMessage, n int) []ChatMessage {
	if len(s) < n {
		return s
	}
	return s[len(s)-n:]
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
