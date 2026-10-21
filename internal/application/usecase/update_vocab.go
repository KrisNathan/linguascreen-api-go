package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
)

type UpdateVocab struct {
	vocabRepo repo.VocabRepo
}

func NewUpdateVocab(vocabRepo repo.VocabRepo) *UpdateVocab {
	return &UpdateVocab{vocabRepo: vocabRepo}
}

func (uc *UpdateVocab) Exec(ctx context.Context, userID, vocabID, translation string) (*aggregate.Vocab, error) {
	v, err := uc.vocabRepo.FindByID(ctx, userID, vocabID)
	if err != nil {
		return nil, err
	}
	if err := v.SetTranslation(translation); err != nil {
		return nil, err
	}
	if err := uc.vocabRepo.UpdateTranslation(ctx, v.ID, v.Translation); err != nil {
		return nil, err
	}
	return v.Vocab, nil
}
