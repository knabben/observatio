package llm

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
)

type Client interface {
	SendMessage(ctx context.Context) (string, error)
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
