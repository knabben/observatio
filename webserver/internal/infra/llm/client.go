package llm

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

type Client interface {
	SendMessage(ctx context.Context, request *ChatMessage, agentID string) (*ChatMessage, error)
	GetClient() anthropic.Client
}

type AnthropicClient struct {
	Client  anthropic.Client
	Service *ObservationService
	Error   string
}

func NewClient() (Client, error) {
	// Create the base client with default settings
	client := &AnthropicClient{
		Client: anthropic.NewClient(),
	}

	// Create the observation service
	service, err := NewObservationService(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create observation service: %w", err)
	}

	// Attach the service to the client
	client.Service = service

	return client, nil
}

func (c *AnthropicClient) GetClient() anthropic.Client {
	return c.Client
}

// SendMessage returns the rendered message to a Websocket or endpoint.
func (c *AnthropicClient) SendMessage(ctx context.Context, request *ChatMessage, agentID string) (*ChatMessage, error) {
	return c.Service.ChatWithAgent(ctx, *request, agentID)
}
