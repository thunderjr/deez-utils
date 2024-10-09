package tool

import (
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
)

type ToolGateway struct {
	tools map[string]Tool
}

func NewToolGateway(tools ...Tool) *ToolGateway {
	t := make(map[string]Tool, len(tools))
	for _, tool := range tools {
		t[tool.Name()] = tool
	}
	return &ToolGateway{t}
}

func (g *ToolGateway) Tool(name string) (Tool, error) {
	tool, ok := g.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool, nil
}

func (g *ToolGateway) Tools() []openai.ChatCompletionToolParam {
	var registeredTools []openai.ChatCompletionToolParam
	for _, tool := range g.tools {
		registeredTools = append(registeredTools, openai.ChatCompletionToolParam{
			Type:     openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(tool.Register()),
		})
	}
	return registeredTools
}

func (g *ToolGateway) HandleToolCalls(calls []openai.ChatCompletionMessageToolCall) ([]openai.ChatCompletionToolMessageParam, error) {
	if len(calls) == 0 {
		return nil, nil
	}

	messages := make([]openai.ChatCompletionToolMessageParam, len(calls))
	for i, toolCall := range calls {
		toolName := toolCall.Function.Name
		tool, err := g.Tool(toolName)
		if err != nil {
			return nil, fmt.Errorf("error getting tool %s: %v", toolName, err)
		}

		var args map[string]any
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			return nil, fmt.Errorf("error unmarshalling tool arguments: %v", err)
		}

		result, err := tool.Execute(args)
		if err != nil {
			return nil, fmt.Errorf("error executing tool %s: %v", toolName, err)
		}

		messages[i] = openai.ToolMessage(toolCall.ID, result)
	}

	return messages, nil
}
