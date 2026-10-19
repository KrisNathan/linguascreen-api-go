package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type GetVocab struct {
	vocabRepo repo.VocabRepo
}

func NewGetVocab(vocabRepo repo.VocabRepo) *GetVocab {
	return &GetVocab{vocabRepo: vocabRepo}
}

func (uc *GetVocab) Exec(ctx context.Context, userID, vocabID string) (*aggregate.Vocab, error) {
	v, err := uc.vocabRepo.FindByID(ctx, userID, vocabID)
	if err != nil {
		return nil, err
	}
	v.Progress = valueobject.Progress{State: valueobject.ProgressNew} // TODO: derive from answered questions once Quiz exists
	return v, nil
}
