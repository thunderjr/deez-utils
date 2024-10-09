package tool

import "github.com/openai/openai-go"

type Tool interface {
	Execute(args map[string]any) (string, error)
	Register() openai.FunctionDefinitionParam
	Name() string
}
