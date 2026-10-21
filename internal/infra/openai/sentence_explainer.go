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
		"{{language_rules}}", explainLanguageRules[learning.String()],
	).Replace(explainSystemPrompt)

	chat, err := e.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: e.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage("Sentence: " + sentence.Text + "\nTranslation: " + sentenceTranslation),
		},
		// Covers reasoning plus a ~300-500 token answer; caps the cost of a rambling reply.
		MaxCompletionTokens: openai.Int(4000),
		ReasoningEffort:     openai.ReasoningEffortLow,
	})
	if err != nil {
		return "", fmt.Errorf("explain sentence: %w", err)
	}
	if len(chat.Choices) == 0 {
		return "", errors.New("explain sentence: no choices")
	}
	choice := chat.Choices[0]
	if choice.FinishReason != "stop" {
		return "", fmt.Errorf("explain sentence: finish reason %q (completion tokens %d, reasoning tokens %d)",
			choice.FinishReason, chat.Usage.CompletionTokens, chat.Usage.CompletionTokensDetails.ReasoningTokens)
	}
	out := strings.TrimSpace(choice.Message.Content)
	if out == "" {
		return "", errors.New("explain sentence: empty output")
	}
	return out, nil
}
