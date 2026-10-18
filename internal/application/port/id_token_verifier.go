package port

import "context"

type IDTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (*Identity, error)
}

type Identity struct {
	Subject string
	Email   string
	Name    string
}
