package aggregate

import (
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

type Question struct {
	ID        string
	SessionID string
	Position  int

	Type         valueobject.QuestionType
	VocabID      string
	OccurrenceID *string // fill_blank only

	Options     []string
	CorrectIdx  int
	AnsweredIdx *int // nil until answered
	Correct     *bool
	AnsweredAt  *time.Time
}

// QuestionDetail is a Question with the vocab and sentence it shows.
type QuestionDetail struct {
	*Question
	Vocab      *Vocab
	Occurrence *Occurrence // fill_blank only
}

func (q *Question) Answer(idx int, now time.Time) error {
	if idx < 0 || idx >= len(q.Options) {
		return errs.WithHint(errs.ErrInvalidInput, "option_index must be between 0 and 3")
	}
	correct := idx == q.CorrectIdx
	q.AnsweredIdx, q.Correct, q.AnsweredAt = &idx, &correct, &now
	return nil
}
