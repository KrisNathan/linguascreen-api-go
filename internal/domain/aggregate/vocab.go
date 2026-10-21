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

const (
	MaxTranslationLength = 200 // runes
	WrongPenalty         = 2
	KnownXP              = 5
)

type Vocab struct {
	ID             string
	CourseID       string
	DictForm       string
	Reading        string
	Translation    string
	XP             int
	CorrectStreak  int
	LastReviewedAt *time.Time // nil until the first answer
	LastSavedAt    time.Time
	CreatedAt      time.Time

	OccurrenceCount int
}

type VocabDetail struct {
	*Vocab
	Occurrences []Occurrence
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

// RecordAnswer updates the practice score, see docs/use_cases/scheduling.md.
func (v *Vocab) RecordAnswer(correct bool, now time.Time) {
	if correct {
		v.XP++
		v.CorrectStreak++
	} else {
		v.XP = max(0, v.XP-WrongPenalty)
		v.CorrectStreak = 0
	}
	v.LastReviewedAt = &now
}

func (v *Vocab) Progress() valueobject.Progress {
	state := valueobject.ProgressLearning
	switch {
	case v.LastReviewedAt == nil:
		state = valueobject.ProgressNew
	case v.XP >= KnownXP:
		state = valueobject.ProgressKnown
	}
	return valueobject.Progress{State: state, Streak: v.CorrectStreak}
}

func normalizeField(s string) string {
	return strings.Join(strings.Fields(norm.NFC.String(s)), " ")
}
