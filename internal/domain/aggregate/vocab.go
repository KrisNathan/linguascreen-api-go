package aggregate

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"golang.org/x/text/unicode/norm"
)

const MaxTranslationLength = 200 // runes

type Vocab struct {
	ID          string
	CourseID    string
	DictForm    string
	Reading     string
	Translation string
	LastSavedAt time.Time
	CreatedAt   time.Time
	Progress    valueobject.Progress

	OccurrenceCount int
	Occurrences     []Occurrence
}

// NewVocab normalizes the fields so the same word from different LLM calls
// maps to the same (course, dict_form, reading) row.
func NewVocab(courseID string, learning valueobject.Language, dictForm, reading, translation string) (*Vocab, error) {
	v := Vocab{
		CourseID: courseID,
		DictForm: normalizeField(dictForm),
		Reading:  normalizeField(reading),
	}
	if learning.String() == "zh" {
		v.Reading = strings.ToLower(v.Reading)
	}
	if v.DictForm == "" {
		return nil, errs.WithHint(errs.ErrInvalidInput, "dict_form is empty")
	}
	if err := v.SetTranslation(translation); err != nil {
		return nil, err
	}
	return &v, nil
}

func (v *Vocab) SetTranslation(translation string) error {
	t := normalizeField(translation)
	if t == "" {
		return errs.WithHint(errs.ErrInvalidInput, "translation is empty")
	}
	if n := utf8.RuneCountInString(t); n > MaxTranslationLength {
		return errs.WithHint(errs.ErrInvalidInput, fmt.Sprintf("translation is %d chars long, max is %d", n, MaxTranslationLength))
	}
	v.Translation = t
	return nil
}

func normalizeField(s string) string {
	return strings.Join(strings.Fields(norm.NFC.String(s)), " ")
}
