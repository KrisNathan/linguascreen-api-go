package repo

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgFlashCardRepo struct {
	pool *pgxpool.Pool
}

func NewPgFlashCardRepo(pool *pgxpool.Pool) *PgFlashCardRepo {
	return &PgFlashCardRepo{pool: pool}
}

func (r *PgFlashCardRepo) SaveReview(ctx context.Context, vocab *aggregate.Vocab, correct bool, reviewedAt time.Time) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO flash_card_reviews (vocab_id, correct, reviewed_at) VALUES ($1, $2, $3)`,
		vocab.ID, correct, reviewedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE vocabs SET xp = $2, correct_streak = $3, last_reviewed_at = $4 WHERE id = $1`,
		vocab.ID, vocab.XP, vocab.CorrectStreak, vocab.LastReviewedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
