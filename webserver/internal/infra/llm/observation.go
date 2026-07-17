package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mcpaggregator "github.com/knabben/observatio/webserver/internal/infra/mcp"
)

const (
	// defaultModel is used when the ANTHROPIC_MODEL env var isn't set. Pinning an explicit
	// constant (with an env override for operators) means upgrading or rolling back the model is
	// a one-line change or a deploy-time config value, never a silent 404 when Anthropic retires
	// a "-latest" style alias out from under a hardcoded string.
	defaultModel anthropic.Model = anthropic.ModelClaudeSonnet5

	// maxToolIterations bounds how many tool-call round trips a single chat turn can make, so a
	// model repeatedly requesting tools can't turn one user message into an unbounded loop.
	maxToolIterations = 5
)

type ObservationService struct {
	// anthropicClient is the client used for interacting with the Anthropic API.
	anthropicClient anthropic.Client

	// model is the Anthropic model used for chat completions.
	model anthropic.Model

	// conversationManager is a map that stores active conversations for each client.
	conversationManager *ConversationManager

	// agents is a map that stores AI agents identified by their unique string IDs.
	agents map[string]*Agent

	// wsConnections maps agent IDs to active WebSocket connections for real-time communication.
	wsConnections map[string]*websocket.Conn

	// aggregator is the shared, process-wide tool source aggregator (built-in kubectl plus any
	// operator-registered external MCP sources) used to render the tool schema offered to Claude
	// and to dispatch tool_use calls. It is constructed once at server startup - not per
	// connection - since building it may involve real MCP handshakes with external sources
	// (specs/009-mcp-server-client-aggregator).
	aggregator *mcpaggregator.Aggregator
}

// NewObservationService creates a per-connection chat service sharing the given process-wide
// Aggregator. aggregator must not be nil.
func NewObservationService(aggregator *mcpaggregator.Aggregator) (*ObservationService, error) {
	service := &ObservationService{
		anthropicClient:     anthropic.NewClient(),
		model:               resolveModel(),
		agents:              make(map[string]*Agent),
		wsConnections:       make(map[string]*websocket.Conn),
		conversationManager: NewConversationManager(5),
		aggregator:          aggregator,
	}
	service.initializeAgents()
	return service, nil
}

// resolveModel lets ANTHROPIC_MODEL override the compiled-in default at deploy time.
func resolveModel() anthropic.Model {
	if m := os.Getenv("ANTHROPIC_MODEL"); m != "" {
		return anthropic.Model(m)
	}
	return defaultModel
}

// StreamChatWithAgent streams the assistant's reply to the given user message: it invokes emit
// with one or more Event:"delta" chunks as text is generated, followed by a final Event:"done"
// chunk once the turn (including any tool calls) has finished. If the model requests a tool call,
// the tool is executed - success or failure - and the result is fed back to the model for a
// follow-up turn, up to maxToolIterations, instead of the tool failure aborting the exchange.
func (s *ObservationService) StreamChatWithAgent(ctx context.Context, message *ChatMessage, emit func(*ChatMessage)) error {
	logger := log.FromContext(ctx)

	userMessage := message.Content
	if s.conversationManager.GetHistoryLength() == 0 {
		userMessage = formatMessage(userMessage)
	}
	s.conversationManager.AddUserMessage(userMessage)

	responseID := generateID()
	history := append([]anthropic.MessageParam{}, s.conversationManager.GetConversationHistory()...)

	var finalText strings.Builder
	var stopReason anthropic.StopReason

	for iteration := 0; iteration < maxToolIterations; iteration++ {
		acc, err := s.streamOnce(ctx, history, responseID, &finalText, emit)
		if err != nil {
			return fmt.Errorf("claude API error: %v", err)
		}
		stopReason = acc.StopReason

		if stopReason != anthropic.StopReasonToolUse {
			break
		}

		history = append(history, acc.ToParam())

		toolResult, ranTools, err := s.runToolCalls(ctx, acc, logger)
		if err != nil {
			return fmt.Errorf("claude API response format error: %v", err)
		}
		if !ranTools {
			break
		}
		history = append(history, toolResult)
	}

	if finalText.Len() == 0 {
		const fallback = "Bot error, try again."
		finalText.WriteString(fallback)
		emit(newStreamChunk(responseID, fallback, "delta"))
	}

	emit(newStreamChunk(responseID, "", "done"))

	logger.Info("Completed response from Claude", "stopReason", stopReason)
	s.conversationManager.AddAssistantMessage(finalText.String())
	s.conversationManager.TrimHistory()
	return nil
}

// streamOnce performs a single streamed request/response turn against the Anthropic API. Text
// deltas are forwarded to emit as they arrive and appended to finalText; the full message is
// accumulated and returned so the caller can inspect its stop reason and any tool_use blocks.
func (s *ObservationService) streamOnce(ctx context.Context, history []anthropic.MessageParam, responseID string, finalText *strings.Builder, emit func(*ChatMessage)) (*anthropic.Message, error) {
	stream := s.anthropicClient.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     s.model,
		MaxTokens: 4000,
		System: []anthropic.TextBlockParam{
			{Text: TASK_SYSTEM},
		},
		Messages: history,
		Tools:    s.aggregator.RenderTools(),
	})

	acc := anthropic.Message{}
	for stream.Next() {
		event := stream.Current()
		if err := acc.Accumulate(event); err != nil {
			return nil, err
		}

		delta, ok := event.AsAny().(anthropic.ContentBlockDeltaEvent)
		if !ok {
			continue
		}
		text, ok := delta.Delta.AsAny().(anthropic.TextDelta)
		if !ok || text.Text == "" {
			continue
		}

		// Strip stray markdown code-fence backticks - a single-character removal is safe to
		// apply per-delta, unlike multi-character tag stripping which would need to buffer
		// across chunks to avoid splitting a tag in half mid-stream.
		chunk := strings.ReplaceAll(text.Text, "`", "")
		finalText.WriteString(chunk)
		emit(newStreamChunk(responseID, chunk, "delta"))
	}
	if err := stream.Err(); err != nil {
		return nil, err
	}

	return &acc, nil
}

// runToolCalls dispatches every tool_use block in the assistant's turn to its owning tool source
// via the Aggregator, and returns a single user MessageParam carrying the matching tool_result
// blocks. Tool-level failures (including "unknown capability" and "source unavailable") are
// reported back to the model as an error result (so it can explain the failure to the operator or
// try something else) rather than aborting the exchange; only a Go-level dispatch error aborts.
func (s *ObservationService) runToolCalls(ctx context.Context, message *anthropic.Message, logger logr.Logger) (anthropic.MessageParam, bool, error) {
	var results []anthropic.ContentBlockParamUnion

	for _, block := range message.Content {
		toolUse, ok := block.AsAny().(anthropic.ToolUseBlock)
		if !ok {
			continue
		}

		output, isError, sourceName, err := s.aggregator.Dispatch(ctx, toolUse.Name, json.RawMessage(toolUse.JSON.Input.Raw()))
		if err != nil {
			return anthropic.MessageParam{}, false, err
		}

		logger.Info("Dispatched tool call", "tool", toolUse.Name, "source", sourceName, "isError", isError)
		results = append(results, anthropic.NewToolResultBlock(toolUse.ID, output, isError))
	}

	if len(results) == 0 {
		return anthropic.MessageParam{}, false, nil
	}
	return anthropic.NewUserMessage(results...), true, nil
}

func generateID() string {
	return uuid.New().String()
}
