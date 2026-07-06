package llm

import (
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

// RunKubectl executes a kubectl command and returns its combined stdout/stderr output. The output
// is returned even when the command fails, so a failure can be reported back to the model as a
// tool result (letting it explain the failure or retry) instead of aborting the whole exchange.
func RunKubectl(command string) (string, error) {
	cmd := exec.Command("kubectl", strings.Fields(command)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
