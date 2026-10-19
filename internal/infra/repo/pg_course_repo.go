package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgCourseRepo struct {
	pool *pgxpool.Pool
}

func NewPgCourseRepo(pool *pgxpool.Pool) *PgCourseRepo {
	return &PgCourseRepo{pool: pool}
}

func (r *PgCourseRepo) Create(ctx context.Context, userID string, baseLanguage, learningLanguage model.Language) (*model.Course, error) {
	c := model.Course{UserID: userID, BaseLanguage: baseLanguage, LearningLanguage: learningLanguage}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO courses (user_id, base_language, learning_language)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, learning_language, base_language) DO UPDATE SET user_id = EXCLUDED.user_id
		 RETURNING id::text, created_at`, userID, baseLanguage.String(), learningLanguage.String()).
		Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// List returns the user's courses, oldest first.
func (r *PgCourseRepo) List(ctx context.Context, userID string) ([]*model.Course, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id::text, user_id::text, base_language, learning_language, created_at
		 FROM courses
		 WHERE user_id = $1
		 ORDER BY created_at`, userID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*model.Course, error) {
		var c model.Course
		var base, learning string
		if err := row.Scan(&c.ID, &c.UserID, &base, &learning, &c.CreatedAt); err != nil {
			return nil, err
		}
		// %v, not %w: an unsupported language in the DB is a server error, not bad client input.
		var err error
		if c.BaseLanguage, err = model.NewLanguage(base); err != nil {
			return nil, fmt.Errorf("course %s: %v", c.ID, err)
		}
		if c.LearningLanguage, err = model.NewLanguage(learning); err != nil {
			return nil, fmt.Errorf("course %s: %v", c.ID, err)
		}
		return &c, nil
	})
}

func (r *PgCourseRepo) FindByID(ctx context.Context, userID, courseID string) (*model.Course, error) {
	var c model.Course
	var base, learning string
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, user_id::text, base_language, learning_language, created_at
		 FROM courses
		 WHERE user_id = $1 AND id = $2`,
		userID, courseID).Scan(&c.ID, &c.UserID, &base, &learning, &c.CreatedAt)
	// 22P02: courseID is not a valid UUID, so no such course exists.
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == "22P02") {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if c.BaseLanguage, err = model.NewLanguage(base); err != nil {
		return nil, fmt.Errorf("course %s: %v", c.ID, err)
	}
	if c.LearningLanguage, err = model.NewLanguage(learning); err != nil {
		return nil, fmt.Errorf("course %s: %v", c.ID, err)
	}
	return &c, nil
}
