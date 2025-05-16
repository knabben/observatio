package llm

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

type Client interface {
	SendMessage(ctx context.Context) (models.LLMResponse, error)
}

type AnthropicClient struct {
	Client  anthropic.Client
	Message string
}

func NewClient(message string) Client {
	return &AnthropicClient{
		Client:  anthropic.NewClient(),
		Message: message,
	}
}
