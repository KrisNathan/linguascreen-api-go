package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

type ListCourse struct {
	courseRepo repo.CourseRepo
}

func NewListCourse(courseRepo repo.CourseRepo) *ListCourse {
	return &ListCourse{courseRepo: courseRepo}
}

func (uc *ListCourse) Exec(ctx context.Context, userID string) ([]*aggregate.Course, error) {
	return uc.courseRepo.List(ctx, userID)
}
