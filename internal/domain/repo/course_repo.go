package repo

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
)

type CourseRepo interface {
	Create(ctx context.Context, userID string, baseLanguage, learningLanguage model.Language) (*model.Course, error)
	List(ctx context.Context, userID string) ([]*model.Course, error)
}
