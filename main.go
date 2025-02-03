package main

import (
	"context"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	llama32  openai.ChatModel = "llama3.2"
	deepseek openai.ChatModel = "deepseek-r1:1.5b"
)

func main() {
	f, err := os.ReadFile("sample-api.json")
	if err != nil {
		panic(err.Error())
	}
	fileContent := string(f)

	chat(fileContent)
}

func chat(userMessage string) string {

	client := openai.NewClient(option.WithBaseURL("http://localhost:11434/v1/"))

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a QA engineer that is responsible for ensuring quality OpenAPI documentation for the company APIs. You are tasked with taking the input JSON and correcting it if it is wrong or summarising what the API does."),
			openai.UserMessage(userMessage),
		}),
		Model: openai.F(llama32),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content
}
