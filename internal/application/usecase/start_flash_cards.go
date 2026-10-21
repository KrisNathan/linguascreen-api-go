package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

const flashCardSize = 20

type StartFlashCards struct {
	courseRepo repo.CourseRepo
	vocabRepo  repo.VocabRepo
}

func NewStartFlashCards(courseRepo repo.CourseRepo, vocabRepo repo.VocabRepo) *StartFlashCards {
	return &StartFlashCards{courseRepo: courseRepo, vocabRepo: vocabRepo}
}

func (uc *StartFlashCards) Exec(ctx context.Context, userID, courseID string) ([]*aggregate.Vocab, error) {
	if _, err := uc.courseRepo.FindByID(ctx, userID, courseID); err != nil {
		return nil, err
	}
	return uc.vocabRepo.ListForPractice(ctx, courseID, flashCardSize)
}
