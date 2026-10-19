package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	errs "github.com/KrisNathan/linguascreen-api-go/internal/errors"
)

type CreateCourse struct {
	courseRepo repo.CourseRepo
}

func NewCreateCourse(courseRepo repo.CourseRepo) *CreateCourse {
	return &CreateCourse{courseRepo: courseRepo}
}

func (uc *CreateCourse) Exec(ctx context.Context, userID, baseLanguage, learningLanguage string) (*model.Course, error) {
	base, err := model.NewLanguage(baseLanguage)
	if err != nil {
		return nil, err
	}

	learning, err := model.NewLanguage(learningLanguage)
	if err != nil {
		return nil, err
	}

	// Compared after normalizing, so "en" and "eng" count as the same language.
	if base == learning {
		return nil, errs.WithHint(errs.ErrInvalidInput, "base_language and learning_language must differ")
	}

	return uc.courseRepo.Create(ctx, userID, base, learning)
}
