package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ObservationService struct {
	// anthropicClient is the client used for interacting with the Anthropic API.
	anthropicClient anthropic.Client

	// conversationManager is a map that stores active conversations for each client.
	conversationManager map[string]ConversationManager

	// agents is a map that stores AI agents identified by their unique string IDs.
	agents map[string]*Agent

	// wsConnections maps agent IDs to active WebSocket connections for real-time communication.
	wsConnections map[string]*websocket.Conn

	// tools represent a collection of tools available for the ObservationService to execute specific operations or commands.
	tools []anthropic.ToolUnionParam
}

func NewObservationService() (*ObservationService, error) {
	service := &ObservationService{
		anthropicClient:     anthropic.NewClient(),
		agents:              make(map[string]*Agent),
		wsConnections:       make(map[string]*websocket.Conn),
		conversationManager: make(map[string]ConversationManager),
		tools:               RenderTools(),
	}
	service.initializeAgents()
	return service, nil
}

func ChatWithAgent(ctx context.Context, message *ChatMessage, agentID string) (*ChatMessage, error) {

	return
}

func (s *ObservationService) ChatWithAgent(ctx context.Context, message *ChatMessage, agentID string) (*ChatMessage, error) {
	logger := log.FromContext(ctx)
	client := s.anthropicClient

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
	s.chatHistory[agentID] = append(s.chatHistory[agentID], []ChatMessage{*message, botMessage}...)
	return &botMessage, nil
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
