package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	llama32          openai.ChatModel = "llama3.2"
	deepseek         openai.ChatModel = "deepseek-r1"
	systemPromptFile string           = "files/system-prompt.txt"
	userMessageFile  string           = "files/sample-api.json"
)

type LLMStructuredResponse struct {
	Status        string `json:"status" jsonschema_description:"Indicates whether the OpenAPI specification is 'compliant' (valid) or 'fixed' (errors corrected)."`
	Message       string `json:"message,omitempty" jsonschema_description:"A brief summary describing any modifications made to fix structural or syntax errors in the OpenAPI specification. This field is omitted if no fixes were needed."`
	CorrectedSpec string `json:"corrected_spec,omitempty" jsonschema_description:"The corrected OpenAPI specification as a JSON object. This field is only present if errors were found and fixed."`
}

var (
	log                            = logrus.New()
	ValidatedOpenAPISchemaResponse = GenerateSchema[LLMStructuredResponse]()
)

func main() {
	ctx := context.Background()
	log.SetFormatter(&logrus.JSONFormatter{})

	jsonFile, error := readFileContents(userMessageFile)
	if error != nil {
		log.Fatal("failed to read JSON file: ", error)
	}

	systemPrompt, error := readFileContents(systemPromptFile)
	if error != nil {
		log.Fatal("failed to read system prompt file: ", error)
	}

	client := openai.NewClient(option.WithBaseURL("http://localhost:11434/v1/"))

	response, err := chat(ctx, client, deepseek, systemPrompt, jsonFile)
	if err != nil {
		log.Fatal("chat request failed: ", err)
	}

	if response.Status == "fixed" {
		// Write response to file
		file, err := os.Create("files/fixed-spec.json")
		if err != nil {
			log.Fatal("failed to create file: ", err)
		}
		defer file.Close()

		_, err = file.WriteString(response.CorrectedSpec)
		if err != nil {
			log.Fatal("failed to write to file: ", err)
		}
	}

	fmt.Println(response.Status, " ", response.Message)
}

func chat(ctx context.Context, client *openai.Client, model openai.ChatModel, systemPrompt, userMessage string) (*LLMStructuredResponse, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("validated open api spec"),
		Description: openai.F("Structured response from LLM for OpenAPI validation and correction"),
		Schema:      openai.F(ValidatedOpenAPISchemaResponse),
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userMessage),
		}),
		Model: openai.F(model),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
	})
	if err != nil {
		return nil, err
	}

	llmResponse := chat.Choices[0].Message.Content
	log.Info(llmResponse)

	var structuredResponse LLMStructuredResponse
	err = json.Unmarshal([]byte(llmResponse), &structuredResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing structured JSON: %w", err)
	}

	return &structuredResponse, nil
}

func readFileContents(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}
