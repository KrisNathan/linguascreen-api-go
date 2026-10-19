package dto

import "github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"

type SavedVocab struct {
	VocabID           string `json:"vocab_id"`
	OccurrenceID      string `json:"occurrence_id"`
	VocabCreated      bool   `json:"vocab_created"`
	OccurrenceCreated bool   `json:"occurrence_created"`
}

func NewSavedVocab(r *usecase.SaveVocabResult) SavedVocab {
	return SavedVocab{
		VocabID:           r.VocabID,
		OccurrenceID:      r.OccurrenceID,
		VocabCreated:      r.VocabCreated,
		OccurrenceCreated: r.OccurrenceCreated,
	}
}
