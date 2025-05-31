package llm

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

var (
	TASK_SYSTEM  = `You will serve as a Kubernetes administrator managing a on-premises datacenter on VMware vCenter.`
	TASK_CONTEXT = `
		Your task is to assist operators in troubleshooting issues within the cluster.
		You should maintain a friendly customer service tone.
		Replace any markdown tags with the appropriate HTML tags.
		Using tooling for troubleshooting is a must, you should be able to identify the root cause of the issue.
		Always ask for the customer before running any tool.
	`
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
	return c.splitResponse(msgContent)
}

func (c *AnthropicClient) splitResponse(response string) (models.LLMResponse, error) {
	descRegex := regexp.MustCompile(`<description>([\s\S]*?)</description>`)
	solRegex := regexp.MustCompile(`<solution>([\s\S]*?)</solution>`)

	response = strings.ReplaceAll(response, "\n", "<br />")
	descMatch := descRegex.FindStringSubmatch(response)
	solMatch := solRegex.FindStringSubmatch(response)

	if len(descMatch) < 2 || len(solMatch) < 2 {
		return models.LLMResponse{}, fmt.Errorf("failed to parse description or solution from response")
	}

	parsed := models.LLMResponse{
		Description: descMatch[1],
		Solution:    solMatch[1],
	}
	return parsed, nil
}

// MessageTemplate defines the structure for formatting error messages
const (
	messageTemplate = "%s\n%s\n%s"
	questionFormat  = "Here is the customer question: <question> %s </question>"
)

// formatMessage creates a formatted message combining the task context, error details,
// and expected output format for the LLM processing
func formatMessage(errorMessage string) string {
	var messageBuilder strings.Builder
	formattedQuestion := fmt.Sprintf(questionFormat, errorMessage)
	messageBuilder.WriteString(fmt.Sprintf(messageTemplate, TASK_CONTEXT, formattedQuestion))

	return messageBuilder.String()
}
