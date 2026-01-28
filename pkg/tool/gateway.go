package tool

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type ToolGateway struct {
	tools map[string]Tool
}

func NewToolGateway(tools ...Tool) *ToolGateway {
	t := make(map[string]Tool, len(tools))
	for _, tool := range tools {
		t[tool.Name()] = tool
	}
	return &ToolGateway{tools: t}
}

func (g *ToolGateway) Tool(name string) (Tool, error) {
	tool, ok := g.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool, nil
}

func (g *ToolGateway) Tools() []llms.FunctionDefinition {
	var registeredTools []llms.FunctionDefinition
	for _, tool := range g.tools {
		registeredTools = append(registeredTools, tool.Register())
	}
	return registeredTools
}

type ToolCall struct {
	ID        string          `json:"id"`
	ToolName  string          `json:"tool"`
	Arguments json.RawMessage `json:"arguments"`
}

func (g *ToolGateway) HandleToolCalls(calls []ToolCall) ([]llms.MessageContent, error) {
	if len(calls) == 0 {
		return nil, nil
	}

	messages := make([]llms.MessageContent, len(calls))
	for i, toolCall := range calls {
		tool, err := g.Tool(toolCall.ToolName)
		if err != nil {
			return nil, fmt.Errorf("error getting tool %s: %v", toolCall.ToolName, err)
		}

		var args map[string]any
		if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
			return nil, fmt.Errorf("error unmarshalling tool arguments: %v", err)
		}

		result, err := tool.Execute(args)
		if err != nil {
			return nil, fmt.Errorf("error executing tool %s: %v", toolCall.ToolName, err)
		}

		messages[i] = llms.TextParts(llms.ChatMessageTypeTool, result)
	}

	return messages, nil
}
