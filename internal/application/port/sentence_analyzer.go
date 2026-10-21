package port

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type SentenceAnalyzer interface {
	Analyze(ctx context.Context, sentence valueobject.Sentence, learning, base valueobject.Language) (*aggregate.Analysis, error)
}
