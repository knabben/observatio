package llm

import "github.com/anthropics/anthropic-sdk-go"

func RenderTools() []anthropic.ToolParam {
	return []anthropic.ToolParam{KubectlTool()}
}

func KubectlTool() anthropic.ToolParam {
	return anthropic.ToolParam{
		Name:        "kubectl",
		Description: anthropic.String("kubectl is a command-line tool for controlling Kubernetes clusters."),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: map[string]interface{}{
				"command": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}
}
