package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

type GetVocab struct {
	vocabRepo repo.VocabRepo
}

func NewGetVocab(vocabRepo repo.VocabRepo) *GetVocab {
	return &GetVocab{vocabRepo: vocabRepo}
}

func (uc *GetVocab) Exec(ctx context.Context, userID, vocabID string) (*aggregate.VocabDetail, error) {
	v, err := uc.vocabRepo.FindByID(ctx, userID, vocabID)
	if err != nil {
		return nil, err
	}
	return v, nil
}
