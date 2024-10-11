package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/thunderjr/deez-utils/pkg/tool"
)

type location struct {
	Name string `json:"name,required"`
}

type calcParams struct {
	Operation string `json:"operation,required"`
	Number1   int    `json:"number1,required"`
	Number2   int    `json:"number2,required"`
}

func getWeather(location *location) (string, error) {
	return "The weather in " + location.Name + " is 72°F and sunny.", nil
}

func calculate(p *calcParams) (string, error) {
	switch p.Operation {
	case "add":
		return fmt.Sprintf("%d + %d = %d", p.Number1, p.Number2, p.Number1+p.Number2), nil
	case "subtract":
		return fmt.Sprintf("%d - %d = %d", p.Number1, p.Number2, p.Number1-p.Number2), nil
	case "multiply":
		return fmt.Sprintf("%d * %d = %d", p.Number1, p.Number2, p.Number1*p.Number2), nil
	case "divide":
		if p.Number2 == 0 {
			return "cannot divide by zero", nil
		}
		return fmt.Sprintf("%d / %d = %d", p.Number1, p.Number2, p.Number1/p.Number2), nil
	default:
		return "unknown operation", nil
	}
}

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "What is the weather in San Francisco?"
	// question := "What is the product of 420 * 69?"

	print("> ")
	println(question)

	// Create a new tool gateway
	toolGateway := tool.NewToolGateway(
		tool.NewTool(
			"get_weather",
			"Get the current weather in a given location",
			getWeather,
		),
		tool.NewTool(
			"calculator",
			"Perform a calculation given a operation (add,subtract,multiply,divide) and two numbers",
			calculate,
		),
	)

	params := openai.ChatCompletionNewParams{
		// Make the tools available to the model
		Tools: openai.F(toolGateway.Tools()),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	// Handle the tool calls
	res, err := toolGateway.HandleToolCalls(completion.Choices[0].Message.ToolCalls)
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		log.Println("Tool response: ", r.Content.String())
	}

	// Call the model with the tool response
	/*
		completion, err = client.Chat.Completions.New(ctx, params)
		if err != nil {
			panic(err)
		}

		log.Println(completion.Choices[0].Message.Content)
	*/

	log.Println("Done")
}
