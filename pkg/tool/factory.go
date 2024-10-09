package tool

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

type toolImpl[Params any] struct {
	fn          func(*Params) (string, error)
	name        string
	description string
}

func NewTool[Params any](name, description string, fn func(*Params) (string, error)) Tool {
	return &toolImpl[Params]{fn, name, description}
}

func (t *toolImpl[Params]) Execute(args map[string]any) (string, error) {
	raw, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("failed to marshal args: %v", err)
	}

	data := new(Params)
	if err := json.Unmarshal(raw, data); err != nil {
		return "", fmt.Errorf("failed to unmarshal args: %v", err)
	}

	return t.fn(data)
}

func (t toolImpl[Params]) Name() string {
	return t.name
}

func (t toolImpl[Params]) Register() shared.FunctionDefinitionParam {
	schema, err := StructToJSONSchema(new(Params))
	if err != nil {
		log.Fatalf("failed to generate JSON schema: %v", err)
	}

	return openai.FunctionDefinitionParam{
		Name:        openai.String(t.Name()),
		Description: openai.String(t.description),
		Parameters:  openai.F(openai.FunctionParameters(schema)),
	}
}
