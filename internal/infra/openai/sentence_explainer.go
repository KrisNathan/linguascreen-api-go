package openai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/openai/openai-go/v3"
)

type SentenceExplainer struct {
	client *openai.Client
	model  string
}

func NewSentenceExplainer(client *openai.Client, model string) *SentenceExplainer {
	return &SentenceExplainer{client: client, model: model}
}

func (e *SentenceExplainer) Explain(ctx context.Context, sentence valueobject.Sentence, sentenceTranslation string, learning, base valueobject.Language) (string, error) {
	prompt := strings.NewReplacer(
		"{{learning_language}}", learning.Name(),
		"{{base_language}}", base.Name(),
	).Replace(explainSystemPrompt)

	chat, err := e.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: e.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage("Sentence: " + sentence.Text + "\nTranslation: " + sentenceTranslation),
		},
		// Expected 300 to 500 tokens; caps the cost of a rambling reply.
		MaxCompletionTokens: openai.Int(1500),
	})
	if err != nil {
		return "", fmt.Errorf("explain sentence: %w", err)
	}
	if len(chat.Choices) == 0 {
		return "", errors.New("explain sentence: no choices")
	}
	choice := chat.Choices[0]
	if choice.FinishReason != "stop" {
		return "", fmt.Errorf("explain sentence: finish reason %q", choice.FinishReason)
	}
	out := strings.TrimSpace(choice.Message.Content)
	if out == "" {
		return "", errors.New("explain sentence: empty output")
	}
	return out, nil
}

// ExplainPromptVersion must be bumped whenever explainSystemPrompt changes,
// so cached results are no longer served.
const (
	ExplainPromptVersion = 1
	// From docs/use_cases/explain.md.
	explainSystemPrompt = `You are a friendly language tutor. The learner is studying {{learning_language}} and wants explanations in {{base_language}}.
They will send a sentence in {{learning_language}} (usually extracted by OCR from a book or a screen) and its translation.

Explain the sentence in {{base_language}}, in Markdown, using only the sections that apply:

## Meaning
What the sentence means here, and how it differs from a word-for-word reading. Stay consistent with the given translation unless it is wrong; if so, say so briefly.

## Grammar
Each grammar point that matters: conjugations, particles, sentence patterns. Name the form (e.g. "て-form + ください"), what it does, and the dictionary form of inflected words.

## Expressions
Idioms or set phrases (e.g. 気をつける), with their meaning as a whole.

## Nuance
Politeness level, tone, register, or anything a learner might misread.

Rules:
- Write readings for words beyond beginner level (Japanese: hiragana in parentheses, Chinese: pinyin with tone marks).
- Be concise: short paragraphs and bullet points, no introduction or closing remarks.
- If the text looks garbled by OCR, say what it most likely says before explaining.`
)
