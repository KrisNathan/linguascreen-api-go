package repo

import (
	"context"
	"errors"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/entity"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPracticeSessionRepo struct {
	pool *pgxpool.Pool
}

func NewPgPracticeSessionRepo(pool *pgxpool.Pool) *PgPracticeSessionRepo {
	return &PgPracticeSessionRepo{pool: pool}
}

func (r *PgPracticeSessionRepo) FindOpen(ctx context.Context, courseID string) (*aggregate.PracticeSession, error) {
	var s aggregate.PracticeSession
	err := r.pool.QueryRow(ctx,
		`
		SELECT id::text, course_id::text, created_at, closed_at
		FROM practice_sessions
		WHERE course_id = $1 AND closed_at IS NULL
		`,
		courseID).Scan(&s.ID, &s.CourseID, &s.CreatedAt, &s.ClosedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if s.Questions, err = r.findQuestions(ctx, s.ID); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PgPracticeSessionRepo) FindByQuestion(ctx context.Context, userID, questionID string) (*aggregate.PracticeSession, error) {
	var s aggregate.PracticeSession
	err := r.pool.QueryRow(ctx,
		`
		SELECT ps.id::text, ps.course_id::text, ps.created_at, ps.closed_at
		FROM questions q
		JOIN practice_sessions ps ON ps.id = q.session_id
		JOIN courses c ON c.id = ps.course_id
		WHERE q.id = $1 AND c.user_id = $2
		`,
		questionID, userID).Scan(&s.ID, &s.CourseID, &s.CreatedAt, &s.ClosedAt)
	// 22P02: questionID is not a valid UUID, so no such question exists.
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == "22P02") {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if s.Questions, err = r.findQuestions(ctx, s.ID); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PgPracticeSessionRepo) findQuestions(ctx context.Context, sessionID string) ([]*aggregate.QuestionDetail, error) {
	rows, err := r.pool.Query(ctx,
		`
		SELECT
			q.id::text,
			q.session_id::text,
			q.position,
			q.type,
			q.vocab_id::text,
			q.occurrence_id::text,
			q.options,
			q.correct_index,
			q.answered_index,
			q.correct,
			q.answered_at,
			v.dict_form,
			v.reading,
			v.translation,
			v.xp,
			v.correct_streak,
			v.last_reviewed_at,
			s.text,
			s.translation,
			o.span_start,
			o.span_end
		FROM questions q
		JOIN vocabs v ON v.id = q.vocab_id
		LEFT JOIN occurrences o ON o.id = q.occurrence_id
		LEFT JOIN sentences s ON s.id = o.sentence_id
		WHERE q.session_id = $1
		ORDER BY q.position
		`,
		sessionID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*aggregate.QuestionDetail, error) {
		q := aggregate.QuestionDetail{Question: &aggregate.Question{}, Vocab: &aggregate.Vocab{}}
		var typ string
		var sentence, sentenceTranslation *string
		var spanStart, spanEnd *int
		err := row.Scan(
			&q.ID, &q.SessionID, &q.Position, &typ, &q.VocabID, &q.OccurrenceID,
			&q.Options, &q.CorrectIdx, &q.AnsweredIdx, &q.Correct, &q.AnsweredAt,
			&q.Vocab.DictForm, &q.Vocab.Reading, &q.Vocab.Translation,
			&q.Vocab.XP, &q.Vocab.CorrectStreak, &q.Vocab.LastReviewedAt,
			&sentence, &sentenceTranslation, &spanStart, &spanEnd,
		)
		if err != nil {
			return nil, err
		}
		q.Type = valueobject.QuestionTypeFromString(typ)
		q.Vocab.ID = q.VocabID
		if q.OccurrenceID != nil {
			q.Occurrence = &entity.Occurrence{
				ID:                  *q.OccurrenceID,
				Sentence:            *sentence,
				SentenceTranslation: *sentenceTranslation,
				Span:                valueobject.Span{Start: *spanStart, End: *spanEnd},
			}
		}
		return &q, nil
	})
}

func (r *PgPracticeSessionRepo) Close(ctx context.Context, sessionID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE practice_sessions SET closed_at = clock_timestamp() WHERE id = $1 AND closed_at IS NULL`,
		sessionID)
	return err
}

func (r *PgPracticeSessionRepo) Create(ctx context.Context, session *aggregate.PracticeSession) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO practice_sessions (course_id) VALUES ($1) RETURNING id::text, created_at`,
		session.CourseID).Scan(&session.ID, &session.CreatedAt)
	// 23505: uk_practice_sessions_open, another start won the race.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return errs.WithHint(errs.ErrConflict, "a quiz is already being started for this course")
	}
	if err != nil {
		return err
	}

	for _, q := range session.Questions {
		q.SessionID = session.ID
		err := tx.QueryRow(ctx,
			`
			INSERT INTO questions (session_id, position, type, vocab_id, occurrence_id, options, correct_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id::text
			`,
			q.SessionID, q.Position, q.Type.String(), q.VocabID, q.OccurrenceID, q.Options, q.CorrectIdx).Scan(&q.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PgPracticeSessionRepo) SaveAnswer(ctx context.Context, q *aggregate.Question, vocab *aggregate.Vocab) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx,
		`UPDATE questions SET answered_index = $2, correct = $3, answered_at = $4 WHERE id = $1 AND answered_at IS NULL`,
		q.ID, q.AnsweredIdx, q.Correct, q.AnsweredAt)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() == 0 {
		return false, nil
	}

	_, err = tx.Exec(ctx,
		`UPDATE vocabs SET xp = $2, correct_streak = $3, last_reviewed_at = $4 WHERE id = $1`,
		vocab.ID, vocab.XP, vocab.CorrectStreak, vocab.LastReviewedAt)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(ctx,
		`
		UPDATE practice_sessions SET closed_at = clock_timestamp()
		WHERE id = $1 AND closed_at IS NULL
			AND NOT EXISTS (SELECT 1 FROM questions WHERE session_id = $1 AND answered_at IS NULL)
		`,
		q.SessionID)
	if err != nil {
		return false, err
	}

	return true, tx.Commit(ctx)
}
