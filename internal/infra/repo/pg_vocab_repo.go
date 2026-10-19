package repo

import (
	"context"
	"errors"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	domainrepo "github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *PgVocabRepo) List(ctx context.Context, courseID string, after *domainrepo.VocabCursor, limit int) ([]*aggregate.Vocab, error) {
	var afterAt *time.Time
	var afterID *string
	if after != nil {
		afterAt, afterID = &after.LastSavedAt, &after.ID
	}
	rows, err := r.pool.Query(ctx,
		`SELECT v.id::text, v.course_id::text, v.dict_form, v.reading, v.translation, v.last_saved_at, v.created_at,
		        (SELECT count(*) FROM occurrences o WHERE o.vocab_id = v.id)
		 FROM vocabs v
		 WHERE v.course_id = $1
		   AND ($2::timestamptz IS NULL OR (v.last_saved_at, v.id) < ($2::timestamptz, $3::uuid))
		 ORDER BY v.last_saved_at DESC, v.id DESC
		 LIMIT $4`,
		courseID, afterAt, afterID, limit)
	if err != nil {
		return nil, err
	}

	vocabs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*aggregate.Vocab, error) {
		var v aggregate.Vocab
		err := row.Scan(&v.ID, &v.CourseID, &v.DictForm, &v.Reading, &v.Translation, &v.LastSavedAt, &v.CreatedAt, &v.OccurrenceCount)
		return &v, err
	})
	// 22P02: the course was already checked, so only the cursor's id can be a malformed UUID.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
		return nil, errs.WithHint(errs.ErrInvalidInput, "invalid cursor")
	}
	return vocabs, err
}

func (r *PgVocabRepo) FindByID(ctx context.Context, userID, vocabID string) (*aggregate.Vocab, error) {
	var v aggregate.Vocab
	err := r.pool.QueryRow(ctx,
		`SELECT v.id::text, v.course_id::text, v.dict_form, v.reading, v.translation, v.last_saved_at, v.created_at
		 FROM vocabs v
		 JOIN courses c ON c.id = v.course_id
		 WHERE v.id = $1 AND c.user_id = $2`,
		vocabID, userID).Scan(&v.ID, &v.CourseID, &v.DictForm, &v.Reading, &v.Translation, &v.LastSavedAt, &v.CreatedAt)
	// 22P02: vocabID is not a valid UUID, so no such vocab exists.
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == "22P02") {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT o.id::text, s.text, s.translation, o.span_start, o.span_end, o.surface, o.translation, o.from_part, o.created_at
		 FROM occurrences o
		 JOIN sentences s ON s.id = o.sentence_id
		 WHERE o.vocab_id = $1
		 ORDER BY o.created_at DESC, o.id DESC`, v.ID)
	if err != nil {
		return nil, err
	}
	v.Occurrences, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (aggregate.Occurrence, error) {
		var o aggregate.Occurrence
		err := row.Scan(&o.ID, &o.Sentence, &o.SentenceTranslation, &o.Span.Start, &o.Span.End, &o.Surface, &o.Translation, &o.FromPart, &o.CreatedAt)
		return o, err
	})
	if err != nil {
		return nil, err
	}
	v.OccurrenceCount = len(v.Occurrences)
	return &v, nil
}
