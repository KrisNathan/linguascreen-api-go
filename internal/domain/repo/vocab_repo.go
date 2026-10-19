package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type VocabRepo interface {
	// Save inserts the sentence, vocab and occurrence in one transaction,
	// reusing any that already exist, and sets vocab.ID and occ.ID.
	Save(ctx context.Context, sentence valueobject.Sentence, sentenceTranslation string, vocab *aggregate.Vocab, occ *aggregate.Occurrence) (vocabCreated, occurrenceCreated bool, err error)
}
