package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/port"
	"github.com/KrisNathan/linguascreen-api-go/internal/application/service"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

type SignInResult struct {
	User   *aggregate.User
	Tokens *service.TokenPair
}

type SignIn struct {
	verifier port.IDTokenVerifier
	userRepo repo.UserRepo
	session  *service.Session
}

func NewSignIn(verifier port.IDTokenVerifier, userRepo repo.UserRepo, session *service.Session) *SignIn {
	return &SignIn{verifier: verifier, userRepo: userRepo, session: session}
}

func (uc *SignIn) Exec(ctx context.Context, idToken string) (*SignInResult, error) {
	identity, err := uc.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}

	user, err := uc.userRepo.Upsert(ctx, identity.Subject, identity.Email, identity.Name)
	if err != nil {
		return nil, err
	}

	tokens, err := uc.session.Issue(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &SignInResult{User: user, Tokens: tokens}, nil
}
