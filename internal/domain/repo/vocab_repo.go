package repo

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/entity"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type VocabCursor struct {
	LastSavedAt time.Time
	ID          string
}

type VocabRepo interface {
	// Save inserts the sentence, vocab and occurrence in one transaction,
	// reusing any that already exist, and sets vocab.ID and occ.ID.
	Save(ctx context.Context, sentence valueobject.Sentence, sentenceTranslation string, vocab *aggregate.Vocab, occ *entity.Occurrence) (vocabCreated, occurrenceCreated bool, err error)

	List(ctx context.Context, courseID string, after *VocabCursor, limit int) ([]*aggregate.Vocab, error)

	FindByID(ctx context.Context, userID, vocabID string) (*aggregate.VocabDetail, error)

	ListForPractice(ctx context.Context, courseID string, limit int) ([]*aggregate.Vocab, error)
	UpdateScore(ctx context.Context, vocabID string, xp, correctStreak int, lastReviewedAt time.Time) error
	UpdateTranslation(ctx context.Context, vocabID, translation string) error

	// ErrNotFound if no vocab owned by user is found
	Delete(ctx context.Context, userID, vocabID string) error

	FindOccurrencesByVocab(ctx context.Context, vocabID string) ([]*entity.Occurrence, error)

	FindDistractorVocabs(ctx context.Context, courseID string, limit int) ([]*aggregate.Vocab, error)

	// FindQuizOccurrence returns the vocab's full-word occurrence quizzed least recently
	// (never quizzed first, ties go to the newest), or nil if it has none.
	FindQuizOccurrence(ctx context.Context, vocabID string) (*entity.Occurrence, error)
}
