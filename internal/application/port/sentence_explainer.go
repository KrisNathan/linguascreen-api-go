package port

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type SentenceExplainer interface {
	Explain(ctx context.Context, sentence valueobject.Sentence, sentenceTranslation string, learning, base valueobject.Language) (string, error)
}
