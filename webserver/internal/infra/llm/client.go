package llm

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

type Client interface {
	SendMessage(ctx context.Context, request string) (models.LLMResponse, error)
	GetClient() anthropic.Client
}

type AnthropicClient struct {
	Client anthropic.Client
	Error  string
}

func NewClient() Client {
	return &AnthropicClient{Client: anthropic.NewClient()}
}

func (c *AnthropicClient) GetClient() anthropic.Client {
	return c.Client
}

func (c *AnthropicClient) SendMessage(ctx context.Context, request string) (response models.LLMResponse, err error) {
	service, err := NewObservationService()
	if err != nil {
		return response, err
	}

	message, err := service.ChatWithAgent(ctx, request)
	if err != nil {
		return response, err
	}

	return c.splitResponse(message.Content)
}
