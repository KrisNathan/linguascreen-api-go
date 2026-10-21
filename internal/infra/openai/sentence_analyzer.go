package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

type SentenceAnalyzer struct {
	client *openai.Client
	model  string
}

func NewSentenceAnalyzer(client *openai.Client, model string) *SentenceAnalyzer {
	return &SentenceAnalyzer{client: client, model: model}
}

func (a *SentenceAnalyzer) Analyze(ctx context.Context, sentence valueobject.Sentence, learning, base valueobject.Language) (*aggregate.Analysis, error) {
	prompt := strings.NewReplacer(
		"{{learning_language}}", learning.Name(),
		"{{base_language}}", base.Name(),
		"{{language_rules}}", languageRules[learning.String()],
	).Replace(systemPrompt)

	chat, err := a.client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Model: a.model,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(prompt),
				openai.UserMessage(sentence.Text),
			},
			Temperature:     openai.Float(0),
			ReasoningEffort: openai.ReasoningEffortLow,
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
					JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:   "sentence_analysis",
						Strict: openai.Bool(true),
						Schema: analysisSchema,
					},
				},
			},
		},
		option.WithJSONSet("provider", map[string]any{"require_parameters": true}),
	)
	if err != nil {
		return nil, fmt.Errorf("analyze sentence: %w", err)
	}
	if len(chat.Choices) == 0 {
		return nil, errors.New("analyze sentence: no choices")
	}
	choice := chat.Choices[0]
	if choice.FinishReason != "stop" {
		return nil, fmt.Errorf("analyze sentence: finish reason %q", choice.FinishReason)
	}

	// Some providers treat the schema as a hint, so the output may still be malformed.
	var out llmAnalysis
	if err := json.Unmarshal([]byte(choice.Message.Content), &out); err != nil {
		return nil, fmt.Errorf("analyze sentence: decode output: %w", err)
	}
	return out.toAggregate(), nil
}

type llmAnalysis struct {
	Translation      string     `json:"translation"`
	LanguageMismatch bool       `json:"language_mismatch"`
	Tokens           []llmToken `json:"tokens"`
}

type llmToken struct {
	Surface     string    `json:"surface"`
	Saveable    bool      `json:"saveable"`
	DictForm    string    `json:"dict_form"`
	Reading     string    `json:"reading"`
	Translation string    `json:"translation"`
	Parts       []llmPart `json:"parts"`
}

type llmPart struct {
	Surface     string `json:"surface"`
	DictForm    string `json:"dict_form"`
	Reading     string `json:"reading"`
	Translation string `json:"translation"`
}

func (o llmAnalysis) toAggregate() *aggregate.Analysis {
	tokens := make([]aggregate.Token, len(o.Tokens))
	for i, t := range o.Tokens {
		parts := make([]aggregate.Part, len(t.Parts))
		for j, p := range t.Parts {
			parts[j] = aggregate.Part{Surface: p.Surface, DictForm: p.DictForm, Reading: p.Reading, Translation: p.Translation}
		}
		tokens[i] = aggregate.Token{
			Surface:     t.Surface,
			Saveable:    t.Saveable,
			DictForm:    t.DictForm,
			Reading:     t.Reading,
			Translation: t.Translation,
			Parts:       parts,
		}
	}
	return &aggregate.Analysis{
		Translation:      o.Translation,
		LanguageMismatch: o.LanguageMismatch,
		Tokens:           tokens,
	}
}

// Strict mode requires every property to be listed in required, so fields that only apply
// to saveable tokens are empty strings otherwise, and parts is an empty array when unused.
var analysisSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"translation":       map[string]any{"type": "string"},
		"language_mismatch": map[string]any{"type": "boolean"},
		"tokens": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"surface":     map[string]any{"type": "string"},
					"saveable":    map[string]any{"type": "boolean"},
					"dict_form":   map[string]any{"type": "string"},
					"reading":     map[string]any{"type": "string"},
					"translation": map[string]any{"type": "string"},
					"parts": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"surface":     map[string]any{"type": "string"},
								"dict_form":   map[string]any{"type": "string"},
								"reading":     map[string]any{"type": "string"},
								"translation": map[string]any{"type": "string"},
							},
							"required":             []string{"surface", "dict_form", "reading", "translation"},
							"additionalProperties": false,
						},
					},
				},
				"required":             []string{"surface", "saveable", "dict_form", "reading", "translation", "parts"},
				"additionalProperties": false,
			},
		},
	},
	"required":             []string{"translation", "language_mismatch", "tokens"},
	"additionalProperties": false,
}
