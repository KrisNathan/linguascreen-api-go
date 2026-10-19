package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgUserRepo struct {
	pool *pgxpool.Pool
}

func NewPgUserRepo(pool *pgxpool.Pool) *PgUserRepo {
	return &PgUserRepo{pool: pool}
}

func (r *PgUserRepo) Upsert(ctx context.Context, googleSub, email, name string) (*aggregate.User, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (google_sub, email, name)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (google_sub) DO UPDATE SET email = EXCLUDED.email, name = EXCLUDED.name
		 RETURNING id::text`, googleSub, email, name).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &aggregate.User{ID: id, GoogleSub: googleSub, Email: email, Name: name}, nil
}
