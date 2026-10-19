package repo

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type VocabCursor struct {
	LastSavedAt time.Time
	ID          string
}

type VocabRepo interface {
	// Save inserts the sentence, vocab and occurrence in one transaction,
	// reusing any that already exist, and sets vocab.ID and occ.ID.
	Save(ctx context.Context, sentence valueobject.Sentence, sentenceTranslation string, vocab *aggregate.Vocab, occ *aggregate.Occurrence) (vocabCreated, occurrenceCreated bool, err error)
	List(ctx context.Context, courseID string, after *VocabCursor, limit int) ([]*aggregate.Vocab, error)
	FindByID(ctx context.Context, userID, vocabID string) (*aggregate.Vocab, error)
	UpdateTranslation(ctx context.Context, vocabID, translation string) error
	// Delete returns errs.ErrNotFound only if the vocab belongs to another user;
	// a missing vocab counts as already deleted.
	Delete(ctx context.Context, userID, vocabID string) error
}
