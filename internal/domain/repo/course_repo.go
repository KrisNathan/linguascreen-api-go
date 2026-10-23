package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/entity"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type CourseRepo interface {
	Create(ctx context.Context, userID string, baseLanguage, learningLanguage valueobject.Language) (*entity.Course, error)
	List(ctx context.Context, userID string) ([]*entity.Course, error)
	FindByID(ctx context.Context, userID, courseID string) (*entity.Course, error)
}
