package dto

import "github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"

type FlashCard struct {
	VocabID     string `json:"vocab_id"`
	DictForm    string `json:"dict_form"`
	Reading     string `json:"reading"`
	Translation string `json:"translation"`
}

type FlashCardDeck struct {
	Cards []FlashCard `json:"cards"`
}

func NewFlashCardDeck(vocabs []*aggregate.Vocab) FlashCardDeck {
	d := FlashCardDeck{Cards: make([]FlashCard, len(vocabs))}
	for i, v := range vocabs {
		d.Cards[i] = FlashCard{VocabID: v.ID, DictForm: v.DictForm, Reading: v.Reading, Translation: v.Translation}
	}
	return d
}
