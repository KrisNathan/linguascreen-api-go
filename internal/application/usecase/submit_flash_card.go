package usecase

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

type SubmitFlashCard struct {
	vocabRepo     repo.VocabRepo
	flashCardRepo repo.FlashCardRepo
}

func NewSubmitFlashCard(vocabRepo repo.VocabRepo, flashCardRepo repo.FlashCardRepo) *SubmitFlashCard {
	return &SubmitFlashCard{vocabRepo: vocabRepo, flashCardRepo: flashCardRepo}
}

func (uc *SubmitFlashCard) Exec(ctx context.Context, userID, vocabID string, correct bool) error {
	v, err := uc.vocabRepo.FindByID(ctx, userID, vocabID)
	if err != nil {
		return err
	}
	now := time.Now()
	v.RecordAnswer(correct, now)
	return uc.flashCardRepo.SaveReview(ctx, v.Vocab, correct, now)
}
