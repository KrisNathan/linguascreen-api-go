package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

type DeleteVocab struct {
	vocabRepo repo.VocabRepo
}

func NewDeleteVocab(vocabRepo repo.VocabRepo) *DeleteVocab {
	return &DeleteVocab{vocabRepo: vocabRepo}
}

func (uc *DeleteVocab) Exec(ctx context.Context, userID, vocabID string) error {
	return uc.vocabRepo.Delete(ctx, userID, vocabID)
}
