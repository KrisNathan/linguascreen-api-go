package google

import (
	"context"
	"errors"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/port"
	"google.golang.org/api/idtoken"
)

type Verifier struct {
	audience  string
	validator *idtoken.Validator
}

func NewVerifier(ctx context.Context, audience string) (*Verifier, error) {
	v, err := idtoken.NewValidator(ctx)
	if err != nil {
		return nil, err
	}
	return &Verifier{audience: audience, validator: v}, nil
}

func (v *Verifier) Verify(ctx context.Context, token string) (*port.Identity, error) {
	p, err := v.validator.Validate(ctx, token, v.audience)
	if err != nil {
		return nil, err
	}

	email, ok := p.Claims["email"].(string)
	if !ok {
		return nil, errors.New("email claim missing")
	}
	name, ok := p.Claims["name"].(string)
	if !ok {
		return nil, errors.New("name claim missing")
	}

	return &port.Identity{Subject: p.Subject, Email: email, Name: name}, nil
}
