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
	jsonFile := fileContents("files/sample-api.json")
	systemPrompt := fileContents("files/system-prompt.md")

	chat(systemPrompt, jsonFile)
}

func chat(systemPrompt, userMessage string) string {
	client := openai.NewClient(option.WithBaseURL("http://localhost:11434/v1/"))

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userMessage),
		}),
		Model: openai.F(llama32),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content
}

func fileContents(path string) string {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	return string(f)
}
