package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ObservationService struct {
	// anthropicClient is the client used for interacting with the Anthropic API.
	anthropicClient anthropic.Client

	// conversationManager is a map that stores active conversations for each client.
	conversationManager *ConversationManager

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
		conversationManager: NewConversationManager(5),
		tools:               RenderTools(),
	}
	service.initializeAgents()
	return service, nil
}

// ChatWithAgent facilitates a conversation between a user and an AI agent by managing message history and API interactions.
func (s *ObservationService) ChatWithAgent(ctx context.Context, message *ChatMessage) (*ChatMessage, error) {
	logger := log.FromContext(ctx)

	userMessage := message.Content
	if s.conversationManager.GetHistoryLength() == 0 {
		userMessage = formatMessage(userMessage)
	}
	s.conversationManager.AddUserMessage(userMessage)

	var historyLength = s.conversationManager.GetHistoryLength()
	if historyLength > 0 {
		logger.Info("Found history messages", "messages", historyLength)
	}

	response, err := s.requestAgent(ctx, s.conversationManager.GetConversationHistory())
	if err != nil {
		return nil, fmt.Errorf("claude API error: %v", err)
	}

	parsedResponse, err := s.responseAgent(response)
	if err != nil {
		return nil, fmt.Errorf("claude API response format error: %v", err)
	}

	if len(response.Content) == 0 {
		return ToMessageParam("Bot error, try again."), nil
	}

	logger.Info("Parsed response from Claude", "response", parsedResponse)
	s.conversationManager.AddAssistantMessage(parsedResponse)
	s.conversationManager.TrimHistory()
	return ToMessageParam(parsedResponse), nil
}

// requestAgent sends a request to the Anthropic API with specified messages and returns the API response or an error.
func (s *ObservationService) requestAgent(ctx context.Context, messages []anthropic.MessageParam) (*anthropic.Message, error) {
	request := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 4000,
		System: []anthropic.TextBlockParam{
			{Text: TASK_SYSTEM},
		},
		Messages: messages,
		Tools:    s.tools,
	}

	return s.anthropicClient.Messages.New(ctx, request)
}

// responseAgent processes the response content from the Anthropic API and constructs a formatted string combining text and tool outputs.
// It handles text blocks and tool-use blocks, extracting detailed outputs as necessary. Returns the formatted response or an error.
func (s *ObservationService) responseAgent(response *anthropic.Message) (string, error) {
	var (
		responseText string
		toolResults  []string
	)

	for _, block := range response.Content {
		switch content := block.AsAny().(type) {
		case anthropic.TextBlock:
			responseText += content.Text

		case anthropic.ToolUseBlock:
			var toolResponse interface{}
			switch block.Name {
			case "kubectl":
				var input struct {
					Command string `json:"command"`
				}
				err := json.Unmarshal([]byte(block.JSON.Input.Raw()), &input)
				if err != nil {
					return "", err
				}
				toolResponse, err = RunKubectl(input.Command)
				if err != nil {
					return "", err
				}
				toolResults = append(toolResults, toolResponse.(string))
			}
		}
	}

	if len(toolResults) > 0 {
		responseText += "\n<tool_results>"
		for _, toolResult := range toolResults {
			responseText += fmt.Sprintf("<pre>%s</pre>", toolResult)
		}
		responseText += "</tool_results>"
	}

	return responseText, nil
}

func generateID() string {
	return uuid.New().String()
}
