package usecase

import (
	"context"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/service"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type StartQuizResult struct {
	SessionID string // empty when there is nothing to practice
	Score     valueobject.QuizScore
	Questions []*aggregate.QuestionDetail
}

type StartQuiz struct {
	courseRepo  repo.CourseRepo
	sessionRepo repo.PracticeSessionRepo
	quiz        *service.QuizService
}

func NewStartQuiz(courseRepo repo.CourseRepo, sessionRepo repo.PracticeSessionRepo, quiz *service.QuizService) *StartQuiz {
	return &StartQuiz{courseRepo: courseRepo, sessionRepo: sessionRepo, quiz: quiz}
}

func (uc *StartQuiz) Exec(ctx context.Context, userID, courseID string) (*StartQuizResult, error) {
	if _, err := uc.courseRepo.FindByID(ctx, userID, courseID); err != nil {
		return nil, err
	}

	openSession, err := uc.sessionRepo.FindOpen(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if openSession != nil {
		if !openSession.Expired(time.Now()) && len(openSession.Unanswered()) > 0 {
			return newStartQuizResult(openSession), nil
		}
		if err := uc.sessionRepo.Close(ctx, openSession.ID); err != nil {
			return nil, err
		}
	}

	session, err := uc.quiz.NewSession(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if len(session.Questions) == 0 {
		return &StartQuizResult{Questions: []*aggregate.QuestionDetail{}}, nil
	}
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return newStartQuizResult(session), nil
}

func newStartQuizResult(s *aggregate.PracticeSession) *StartQuizResult {
	return &StartQuizResult{SessionID: s.ID, Score: s.Score(), Questions: s.Unanswered()}
}
