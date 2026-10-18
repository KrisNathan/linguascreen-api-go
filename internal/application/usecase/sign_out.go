package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/util"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

type SignOut struct {
	refreshTokenRepo repo.RefreshTokenRepo
}

func NewSignOut(refreshTokenRepo repo.RefreshTokenRepo) *SignOut {
	return &SignOut{refreshTokenRepo: refreshTokenRepo}
}

// Exec revokes the refresh token. Unknown or already revoked tokens are not an error.
func (uc *SignOut) Exec(ctx context.Context, refreshToken string) error {
	return uc.refreshTokenRepo.RevokeToken(ctx, util.HashRefreshToken(refreshToken))
}
