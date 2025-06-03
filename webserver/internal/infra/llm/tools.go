package llm

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

func RenderTools() []anthropic.ToolUnionParam {
	allTools := []anthropic.ToolParam{KubectlTool()}
	tools := make([]anthropic.ToolUnionParam, len(allTools))
	for i, toolParam := range allTools {
		tools[i] = anthropic.ToolUnionParam{OfTool: &toolParam}
	}
	return tools
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

func RunKubectl(command string) (string, error) {
	// Execute kubectl command using os/exec
	cmd := exec.Command("kubectl", strings.Fields(command)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing kubectl command: %v", err)
	}

	return string(output), nil
}
