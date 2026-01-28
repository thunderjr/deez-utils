package tool

import "github.com/tmc/langchaingo/llms"

type Tool interface {
	Execute(args map[string]any) (string, error)
	Name() string
	Register() llms.FunctionDefinition
}
