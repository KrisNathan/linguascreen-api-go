package usecase

import (
	"context"
	"errors"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/service"
	"github.com/KrisNathan/linguascreen-api-go/internal/application/util"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	errs "github.com/KrisNathan/linguascreen-api-go/internal/errors"
)

type RefreshToken struct {
	refreshTokenRepo repo.RefreshTokenRepo
	session          *service.Session
}

func NewRefreshToken(refreshTokenRepo repo.RefreshTokenRepo, session *service.Session) *RefreshToken {
	return &RefreshToken{refreshTokenRepo: refreshTokenRepo, session: session}
}

func (uc *RefreshToken) Exec(ctx context.Context, refreshToken string) (*service.TokenPair, error) {
	oldHash := util.HashRefreshToken(refreshToken)

	userID, err := uc.refreshTokenRepo.Consume(ctx, oldHash)
	if errors.Is(err, errs.ErrNotFound) {
		// Missing, expired or already revoked. A revoked token being reused means it was likely stolen.
		if t, findErr := uc.refreshTokenRepo.FindByHash(ctx, oldHash); findErr == nil && t.RevokedAt != nil {
			if err := uc.refreshTokenRepo.RevokeAllForUser(ctx, t.UserID); err != nil {
				return nil, err
			}
		}
		return nil, errs.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	return uc.session.Issue(ctx, userID)
}
