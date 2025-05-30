package llm

type ToolSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required"`
}

type Tool struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	InputSchema ToolSchema `json:"input_schema"`
}
