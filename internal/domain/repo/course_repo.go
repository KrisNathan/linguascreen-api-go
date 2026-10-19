package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type CourseRepo interface {
	Create(ctx context.Context, userID string, baseLanguage, learningLanguage valueobject.Language) (*aggregate.Course, error)
	List(ctx context.Context, userID string) ([]*aggregate.Course, error)
	FindByID(ctx context.Context, userID, courseID string) (*aggregate.Course, error)
}
