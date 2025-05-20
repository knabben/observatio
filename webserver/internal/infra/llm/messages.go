package llm

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

var (
	TASK_SYSTEM = `
		You will serve as a Kubernetes administrator managing an on-premises datacenter using VMware vCenter.
	`
	TASK_CONTEXT = `
		Your task is to assist operators in troubleshooting issues within the cluster.
		You should maintain a friendly customer service tone.
	`
	INPUT_DATA  = `Here is the customer question: <question> %s </question>`
	OUTPUT_DATA = `The answer must be divided in two parts:
		A verbose description of the error in <description></description> tags
		The solution of the error in <solution></solution> tags`
)

func (c *AnthropicClient) formatMessage(msg string) string {
	data := fmt.Sprintf(INPUT_DATA, msg)
	return fmt.Sprintf("%s\n%s\n%s", TASK_CONTEXT, data, OUTPUT_DATA)
}

func (c *AnthropicClient) SendMessage(ctx context.Context) (response models.LLMResponse, err error) {
	var msg *anthropic.Message

	msg, err = c.Client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(c.formatMessage(c.Error))),
		},
		System: []anthropic.TextBlockParam{
			{Text: TASK_SYSTEM},
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
