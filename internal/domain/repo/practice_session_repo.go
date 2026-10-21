package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
)

type PracticeSessionRepo interface {
	// FindOpen returns nil if the course has no open session.
	FindOpen(ctx context.Context, courseID string) (*aggregate.PracticeSession, error)

	// ErrNotFound if the question doesn't exist or belongs to another user.
	FindByQuestion(ctx context.Context, userID, questionID string) (*aggregate.PracticeSession, error)

	Close(ctx context.Context, sessionID string) error

	// ErrConflict if the course already has an open session.
	Create(ctx context.Context, session *aggregate.PracticeSession) error

	// SaveAnswer also saves the vocab's score and closes the session after its last answer.
	// It saves nothing and returns false if the question was already answered.
	SaveAnswer(ctx context.Context, q *aggregate.Question, vocab *aggregate.Vocab) (bool, error)
}
