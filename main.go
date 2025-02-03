package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	llama32      openai.ChatModel = "llama3.2"
	deepseek     openai.ChatModel = "deepseek-r1:1.5b"
	systemPrompt string           = "files/system-prompt.md"
	userMessage  string           = "files/sample-api.json"
)

func main() {
	ctx := context.Background()

	jsonFile, error := readFileContents("files/sample-api.json")
	if error != nil {
		log.Fatalf("failed to read JSON file: %v", error)
	}

	systemPrompt, error := readFileContents("files/system-prompt.md")
	if error != nil {
		log.Fatalf("failed to read system prompt file: %v", error)
	}

	client := openai.NewClient(option.WithBaseURL("http://localhost:11434/v1/"))

	response, err := chat(ctx, client, llama32, systemPrompt, jsonFile)
	if err != nil {
		log.Fatalf("chat request failed: %v", err)
	}

	fmt.Println("Chat response:", response)
}

func chat(ctx context.Context, client *openai.Client, model openai.ChatModel, systemPrompt, userMessage string) (string, error) {
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userMessage),
		}),
		Model: openai.F(model),
	})
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func readFileContents(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}
