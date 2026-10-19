package port

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type SentenceAnalyzer interface {
	Analyze(ctx context.Context, sentence valueobject.Sentence, learning, base model.Language) (*aggregate.Analysis, error)
}
