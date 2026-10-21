package dto

import (
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type Analysis struct {
	Sentence         string  `json:"sentence"`
	Translation      string  `json:"translation"`
	LanguageMismatch bool    `json:"language_mismatch"`
	Tokens           []Token `json:"tokens"`
}

type Token struct {
	Surface     string `json:"surface"`
	Span        [2]int `json:"span"`
	Saveable    bool   `json:"saveable"`
	DictForm    string `json:"dict_form,omitempty"`
	Reading     string `json:"reading,omitempty"`
	Translation string `json:"translation,omitempty"`
	Parts       []Part `json:"parts,omitempty"`
}

type Part struct {
	Surface     string `json:"surface"`
	Span        [2]int `json:"span"`
	DictForm    string `json:"dict_form"`
	Reading     string `json:"reading"`
	Translation string `json:"translation"`
}

func NewAnalysis(s *valueobject.Sentence, a *aggregate.Analysis) Analysis {
	tokens := make([]Token, len(a.Tokens))
	for i, t := range a.Tokens {
		tokens[i] = Token{
			Surface:     t.Surface,
			Span:        [2]int{t.Span.Start, t.Span.End},
			Saveable:    t.Saveable,
			DictForm:    t.DictForm,
			Reading:     t.Reading,
			Translation: t.Translation,
		}
		for _, p := range t.Parts {
			tokens[i].Parts = append(tokens[i].Parts, Part{
				Surface:     p.Surface,
				Span:        [2]int{p.Span.Start, p.Span.End},
				DictForm:    p.DictForm,
				Reading:     p.Reading,
				Translation: p.Translation,
			})
		}
	}
	return Analysis{
		Sentence:         s.Text,
		Translation:      a.Translation,
		LanguageMismatch: a.LanguageMismatch,
		Tokens:           tokens,
	}
}
