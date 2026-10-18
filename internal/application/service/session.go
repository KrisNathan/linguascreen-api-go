package service

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/port"
	"github.com/KrisNathan/linguascreen-api-go/internal/application/util"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

const refreshTTL = 30 * 24 * time.Hour

type TokenPair struct {
	AccessToken     string
	AccessExpiresIn time.Duration
	RefreshToken    string
}

type Session struct {
	refreshTokenRepo repo.RefreshTokenRepo
	tokenService     port.TokenService
}

func NewSession(refreshTokenRepo repo.RefreshTokenRepo, tokenService port.TokenService) *Session {
	return &Session{refreshTokenRepo: refreshTokenRepo, tokenService: tokenService}
}

// Issue stores a new refresh token for userID and returns it with a fresh access token.
func (s *Session) Issue(ctx context.Context, userID string) (*TokenPair, error) {
	refreshToken, refreshHash, err := util.NewRefreshToken()
	if err != nil {
		return nil, err
	}
	if err := s.refreshTokenRepo.Create(ctx, userID, refreshHash, time.Now().Add(refreshTTL)); err != nil {
		return nil, err
	}

	access, err := s.tokenService.IssueAccess(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:     access,
		AccessExpiresIn: s.tokenService.AccessTTL(),
		RefreshToken:    refreshToken,
	}, nil
}
