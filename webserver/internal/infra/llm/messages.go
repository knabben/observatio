package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

var (
	TASK_SYSTEM  = `You will serve as a Kubernetes administrator managing a on-premises datacenter on VMware vCenter.`
	TASK_CONTEXT = `Your task is to assist operators in troubleshooting issues within the cluster.
Provide a detailed explanation of the issue. New inputs are provided and you must respond accordingly the context
Do not use any existent tool for now.`
)

func (c *AnthropicClient) SendMessageMove(ctx context.Context) (response models.LLMResponse, err error) {
	var msg *anthropic.Message

	msg, err = c.Client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(formatMessage(c.Error))),
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
	return models.LLMResponse{Result: msgContent}, nil
}

// MessageTemplate defines the structure for formatting error messages
const (
	messageTemplate = "%s\n%s"
	questionFormat  = "Here is the customer question: <question>%s</question>"
)

// formatMessage creates a formatted message combining the task context, error details,
// and expected output format for the LLM processing
func formatMessage(errorMessage string) string {
	var messageBuilder strings.Builder
	formattedQuestion := fmt.Sprintf(questionFormat, errorMessage)
	messageBuilder.WriteString(fmt.Sprintf(messageTemplate, TASK_CONTEXT, formattedQuestion))
	return messageBuilder.String()
}
