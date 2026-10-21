package aggregate

import (
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

const SessionTTL = 24 * time.Hour

type PracticeSession struct {
	ID       string
	CourseID string

	CreatedAt time.Time
	ClosedAt  *time.Time

	Questions []*QuestionDetail
}

func (s *PracticeSession) Expired(now time.Time) bool {
	return s.ClosedAt != nil || now.Sub(s.CreatedAt) >= SessionTTL
}

func (s *PracticeSession) Unanswered() []*QuestionDetail {
	out := []*QuestionDetail{}
	for _, q := range s.Questions {
		if q.AnsweredAt == nil {
			out = append(out, q)
		}
	}
	return out
}

func (s *PracticeSession) Question(id string) *QuestionDetail {
	for _, q := range s.Questions {
		if q.ID == id {
			return q
		}
	}
	return nil
}

func (s *PracticeSession) Score() valueobject.QuizScore {
	score := valueobject.QuizScore{Total: len(s.Questions)}
	for _, q := range s.Questions {
		if q.AnsweredAt != nil {
			score.Answered++
			if *q.Correct {
				score.Correct++
			}
		}
	}
	return score
}
