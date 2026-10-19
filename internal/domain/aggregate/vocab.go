package aggregate

import (
	"strings"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"golang.org/x/text/unicode/norm"
)

type Vocab struct {
	ID          string
	CourseID    string
	DictForm    string
	Reading     string
	Translation string
}

// NewVocab normalizes the fields so the same word from different LLM calls
// maps to the same (course, dict_form, reading) row.
func NewVocab(courseID string, learning valueobject.Language, dictForm, reading, translation string) (*Vocab, error) {
	v := Vocab{
		CourseID:    courseID,
		DictForm:    normalizeField(dictForm),
		Reading:     normalizeField(reading),
		Translation: normalizeField(translation),
	}
	if learning.String() == "zh" {
		v.Reading = strings.ToLower(v.Reading)
	}
	if v.DictForm == "" {
		return nil, errs.WithHint(errs.ErrInvalidInput, "dict_form is empty")
	}
	if v.Translation == "" {
		return nil, errs.WithHint(errs.ErrInvalidInput, "translation is empty")
	}
	return &v, nil
}

func normalizeField(s string) string {
	return strings.Join(strings.Fields(norm.NFC.String(s)), " ")
}
