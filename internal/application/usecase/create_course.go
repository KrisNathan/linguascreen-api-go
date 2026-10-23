package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/entity"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

type CreateCourse struct {
	courseRepo repo.CourseRepo
}

func NewCreateCourse(courseRepo repo.CourseRepo) *CreateCourse {
	return &CreateCourse{courseRepo: courseRepo}
}

func (uc *CreateCourse) Exec(ctx context.Context, userID, baseLanguage, learningLanguage string) (*entity.Course, error) {
	base, err := valueobject.NewLanguage(baseLanguage)
	if err != nil {
		return nil, err
	}

	learning, err := valueobject.NewLanguage(learningLanguage)
	if err != nil {
		return nil, err
	}

	if base == learning {
		return nil, errs.WithHint(errs.ErrInvalidInput, "base_language and learning_language must differ")
	}

	return uc.courseRepo.Create(ctx, userID, base, learning)
}
