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
			anthropic.NewUserMessage(anthropic.NewTextBlock(c.Error)),
		},
		Model:         anthropic.ModelClaude3_7SonnetLatest,
		StopSequences: []string{"---\n"},
	})
	if err != nil {
		return response, err
	}
	var msgContent string
	for _, m := range msg.Content {
		msgContent += m.Text
	}
	return models.LLMResponse{Data: msgContent + msg.StopSequence}, nil
}
