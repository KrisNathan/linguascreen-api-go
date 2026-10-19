package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
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

func (a *SentenceAnalyzer) Analyze(ctx context.Context, sentence valueobject.Sentence, learning, base model.Language) (*aggregate.Analysis, error) {
	prompt := strings.NewReplacer(
		"{{learning_language}}", learning.Name(),
		"{{base_language}}", base.Name(),
	).Replace(systemPrompt)

	chat, err := a.client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Model: a.model,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(prompt),
				openai.UserMessage(sentence.Text),
			},
			Temperature: openai.Float(0),
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

// PromptVersion must be bumped whenever systemPrompt or analysisSchema changes,
// so cached results are no longer served.
const (
	PromptVersion = 1

	// From docs/use_cases/analyze_sentence.md, adjusted for strict mode (empty strings and arrays instead of omitted fields).
	systemPrompt = `You are a language-learning assistant. The user is learning {{learning_language}} and their native language is {{base_language}}. They will send you one sentence in {{learning_language}}, usually extracted by OCR from a book or a screen.

Your job:
1. Translate the whole sentence into {{base_language}}.
2. Split the sentence into tokens so the learner can tap any word to see what it means.
3. Set language_mismatch to true if the sentence is mostly not in {{learning_language}}. Still translate and tokenize it.

Tokens:
- Copy the sentence exactly. The surfaces of all tokens, joined in order, must reproduce the sentence character for character. Do not fix typos or OCR errors, do not convert characters, punctuation or width (keep "！" as "！"), do not add or drop spaces. Punctuation, symbols and spaces are tokens too.
- A token is a unit a learner would look up in a dictionary. Keep an inflected word together with its inflection and auxiliaries as one token (e.g. 読みます, 走った, 食べさせられなかった), not split into grammatical pieces.
- saveable is false for punctuation, symbols, spaces, numbers, and grammatical function words that carry no dictionary meaning on their own (e.g. Japanese particles like を, が, に, Chinese 的, 了). Otherwise true. Only saveable tokens need dict_form, reading and translation; for other tokens set them to empty strings and parts to an empty array.

For each saveable token:
- dict_form: the dictionary (lemma) form, e.g. 読みます → 読む, 走った → 走る. For languages that do not inflect, the word itself.
- reading: the pronunciation of dict_form.
  - Japanese: hiragana. For words written only in kana, repeat dict_form.
  - Chinese: pinyin with tone marks, syllables separated by spaces (e.g. yín háng). Pick the pronunciation that fits this sentence (e.g. 长 is cháng in 很长, zhǎng in 长大).
  - Languages without a separate reading system: empty string.
- translation: a short dictionary-style meaning in {{base_language}} for the sense used in this sentence, 1 to 4 words, no explanations. In English, write verbs as "to ..." (e.g. "to read").

Parts:
- If a token is a compound made of smaller words that are worth learning on their own, add parts (e.g. 火车站 → 火车 + 站, 国際連合 → 国際 + 連合). Otherwise set parts to an empty array.
- Do not split inflections or grammatical endings into parts.
- The part surfaces, joined in order, must reproduce the token's surface exactly.
- Each part has its own dict_form, reading and translation, following the rules above. Parts have no parts.

Sentence translation:
- Natural and faithful, in {{base_language}}. If the text looks garbled by OCR, translate your best interpretation, but still copy the token surfaces exactly as given.`
)
