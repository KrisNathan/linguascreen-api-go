package repo

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
)

type FlashCardRepo interface {
	// SaveReview also saves the vocab's score.
	SaveReview(ctx context.Context, vocab *aggregate.Vocab, correct bool, reviewedAt time.Time) error
}
