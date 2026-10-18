package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
)

type UserRepo interface {
	// Upsert creates the user for googleSub, or updates its email and name if it already exists.
	Upsert(ctx context.Context, googleSub, email, name string) (*model.User, error)
}
