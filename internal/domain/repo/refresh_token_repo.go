package repo

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
)

type RefreshTokenRepo interface {
	Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	// Consume atomically revokes an active (unrevoked, unexpired) token and returns its user ID.
	// Returns errs.ErrNotFound if there is no such active token.
	Consume(ctx context.Context, tokenHash string) (userID string, err error)
	// FindByHash returns the token in any state. Returns errs.ErrNotFound if it does not exist.
	FindByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	RevokeToken(ctx context.Context, tokenHash string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
