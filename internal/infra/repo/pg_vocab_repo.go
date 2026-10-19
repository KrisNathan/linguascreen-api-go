package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgVocabRepo struct {
	pool *pgxpool.Pool
}

func NewPgVocabRepo(pool *pgxpool.Pool) *PgVocabRepo {
	return &PgVocabRepo{pool: pool}
}

// The no-op DO UPDATEs make RETURNING give back existing rows; xmax = 0 only
// for a row this statement inserted. Existing rows keep their translation.
func (r *PgVocabRepo) Save(ctx context.Context, sentence valueobject.Sentence, sentenceTranslation string, vocab *aggregate.Vocab, occ *aggregate.Occurrence) (vocabCreated, occurrenceCreated bool, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, false, err
	}
	defer tx.Rollback(ctx)

	var sentenceID string
	err = tx.QueryRow(ctx,
		`INSERT INTO sentences (course_id, text, text_hash, translation)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (course_id, text_hash) DO UPDATE SET course_id = EXCLUDED.course_id
		 RETURNING id::text`,
		vocab.CourseID, sentence.Text, sentence.Hash(), sentenceTranslation).Scan(&sentenceID)
	if err != nil {
		return false, false, err
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO vocabs (course_id, dict_form, reading, translation)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (course_id, dict_form, reading) DO UPDATE SET course_id = EXCLUDED.course_id
		 RETURNING id::text, xmax = 0`,
		vocab.CourseID, vocab.DictForm, vocab.Reading, vocab.Translation).Scan(&vocab.ID, &vocabCreated)
	if err != nil {
		return false, false, err
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO occurrences (vocab_id, sentence_id, span_start, span_end, surface, translation, from_part)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (vocab_id, sentence_id, span_start) DO UPDATE SET vocab_id = EXCLUDED.vocab_id
		 RETURNING id::text, xmax = 0`,
		vocab.ID, sentenceID, occ.Span.Start, occ.Span.End, occ.Surface, occ.Translation, occ.FromPart).
		Scan(&occ.ID, &occurrenceCreated)
	if err != nil {
		return false, false, err
	}

	// A new vocab already has last_saved_at = now from the default.
	if occurrenceCreated && !vocabCreated {
		if _, err := tx.Exec(ctx, `UPDATE vocabs SET last_saved_at = clock_timestamp() WHERE id = $1`, vocab.ID); err != nil {
			return false, false, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, false, err
	}
	return vocabCreated, occurrenceCreated, nil
}
