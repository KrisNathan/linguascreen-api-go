package services

import (
	"context"
	"encoding/json"
	"linguascreen/models"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

type OpenAILLMService struct {
	client openai.Client
}

func NewLLMService() *OpenAILLMService {
	return &OpenAILLMService{
		client: openai.NewClient(
			option.WithBaseURL(os.Getenv("OPENAI_BASE_URL")),
			option.WithAPIKey(os.Getenv("OPENAI_API_KEY")), // or set OPENAI_API_KEY in your env
		),
	}
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

// Generate the JSON schema at initialization time
var ExplainResponseSchema = GenerateSchema[models.ExplainResult]()

func (s *OpenAILLMService) Explain(sentence string, targetLang string, userLang string, nmtTranslation string, memories []models.Memory) (*models.ExplainResult, error) {

	memoriesJsonBytes, err := json.Marshal(memories)
	if err != nil {
		return nil, err
	}
	memoriesJson := string(memoriesJsonBytes)

	system := "You are a helpful assistant that provides structured explanations of new knowledge in language learning. When given a sentence, identify new words, phrases, or grammar points that a learner might not know. For each item, provide a detailed explanation, translation, and an example sentence. Format your response strictly according to the provided JSON schema. The user is learning " + targetLang + ". The user knows " + userLang + ". If the sentence is in " + targetLang + ", respond in " + userLang + ". If the sentence is in " + userLang + ", translate the sentence to " + targetLang + " and explain it in " + userLang + "." + "\nUser memories:\n" + memoriesJson

	question := "Explain the following text and identify new words, phrases, or grammar points that a learner might not know: " + sentence

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "explain_response",
		Description: openai.String("Structured response containing new knowledge explanations"),
		Schema:      ExplainResponseSchema,
		Strict:      openai.Bool(true),
	}

	chat, chatErr := s.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage(question),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		// only certain models can perform structured outputs
		Model: openai.ChatModelGPT4o2024_08_06,
	})
	if chatErr != nil {
		return nil, chatErr
	}

	// extract into a well-typed struct
	var explainResponse models.ExplainResult
	unmarshallErr := json.Unmarshal([]byte(chat.Choices[0].Message.Content), &explainResponse)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}

	return &explainResponse, nil
}

var TextRearrangeResponseSchema = GenerateSchema[models.TextRearrangeResult]()

func (s *OpenAILLMService) RearrangeText(text string) (string, error) {

	system := "You are a helpful assistant that rearranges words in a sentence to improve clarity. Minimize changes to the original text. Avoid changing the meaning. Preserve punctuation."
	user := "Rearrange the following text for better clarity: " + text

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "sentence_rearrangement",
		Description: openai.String("Structured response containing the rearranged text"),
		Schema:      TextRearrangeResponseSchema,
		Strict:      openai.Bool(true),
	}

	chat, chatErr := s.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage(user),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		Model: openai.ChatModelGPT4o2024_08_06,
	})
	if chatErr != nil {
		return "", chatErr
	}

	// extract into a well-typed struct
	var rearrangeResponse models.TextRearrangeResult
	unmarshallErr := json.Unmarshal([]byte(chat.Choices[0].Message.Content), &rearrangeResponse)
	if unmarshallErr != nil {
		return "", unmarshallErr
	}

	return rearrangeResponse.ProcessedText, nil
}
