package usecase

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

type SubmitQuizResult struct {
	Question         *aggregate.QuestionDetail
	Score            valueobject.QuizScore
	SessionCompleted bool
}

type SubmitQuiz struct {
	sessionRepo repo.PracticeSessionRepo
}

func NewSubmitQuiz(sessionRepo repo.PracticeSessionRepo) *SubmitQuiz {
	return &SubmitQuiz{sessionRepo: sessionRepo}
}

func (uc *SubmitQuiz) Exec(ctx context.Context, userID, questionID string, optionIdx int) (*SubmitQuizResult, error) {
	session, err := uc.sessionRepo.FindByQuestion(ctx, userID, questionID)
	if err != nil {
		return nil, err
	}
	q := session.Question(questionID)

	// Checked before expiry so a retry gets the same result even after the session closed.
	if q.AnsweredAt != nil {
		return newSubmitQuizResult(session, q), nil
	}

	now := time.Now()
	if session.Expired(now) {
		return nil, errs.WithHint(errs.ErrConflict, "quiz session has expired")
	}
	if err := q.Answer(optionIdx, now); err != nil {
		return nil, err
	}
	q.Vocab.RecordAnswer(*q.Correct, now)

	saved, err := uc.sessionRepo.SaveAnswer(ctx, q.Question, q.Vocab)
	if err != nil {
		return nil, err
	}
	if !saved {
		// Lost a race with another submit of the same question, return what it stored.
		if session, err = uc.sessionRepo.FindByQuestion(ctx, userID, questionID); err != nil {
			return nil, err
		}
		q = session.Question(questionID)
	}
	return newSubmitQuizResult(session, q), nil
}

func newSubmitQuizResult(s *aggregate.PracticeSession, q *aggregate.QuestionDetail) *SubmitQuizResult {
	return &SubmitQuizResult{Question: q, Score: s.Score(), SessionCompleted: len(s.Unanswered()) == 0}
}
