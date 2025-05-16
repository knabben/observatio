package llm

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

func (c *AnthropicClient) SendMessage(ctx context.Context) (response models.LLMResponse, err error) {
	var msg *anthropic.Message
	msg, err = c.Client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(c.Message)),
		},
		Model:         anthropic.ModelClaude3_7SonnetLatest,
		StopSequences: []string{"```\n"},
	})
	if err != nil {
		return response, err
	}
	return models.LLMResponse{Response: msg.Content[0].Text + msg.StopSequence}, nil
}
