package repo

import (
	"context"
	"errors"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRefreshTokenRepo struct {
	pool *pgxpool.Pool
}

func NewPgRefreshTokenRepo(pool *pgxpool.Pool) *PgRefreshTokenRepo {
	return &PgRefreshTokenRepo{pool: pool}
}

func (r *PgRefreshTokenRepo) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1::uuid, $2, $3)`, userID, tokenHash, expiresAt)
	return err
}

func (r *PgRefreshTokenRepo) Consume(ctx context.Context, tokenHash string) (string, error) {
	var userID string
	err := r.pool.QueryRow(ctx,
		`UPDATE refresh_tokens SET revoked_at = now()
		 WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > now()
		 RETURNING user_id::text`, tokenHash).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errs.ErrNotFound
	}
	return userID, err
}

func (r *PgRefreshTokenRepo) FindByHash(ctx context.Context, tokenHash string) (*aggregate.RefreshToken, error) {
	var t aggregate.RefreshToken
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, user_id::text, expires_at, revoked_at
		 FROM refresh_tokens WHERE token_hash = $1`, tokenHash).
		Scan(&t.ID, &t.UserID, &t.ExpiresAt, &t.RevokedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PgRefreshTokenRepo) RevokeToken(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = now()
		 WHERE token_hash = $1 AND revoked_at IS NULL`, tokenHash)
	return err
}

func (r *PgRefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = now()
		 WHERE user_id = $1::uuid AND revoked_at IS NULL`, userID)
	return err
}
