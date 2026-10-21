package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/port"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

type ExplainSentence struct {
	courseRepo repo.CourseRepo
	explainer  port.SentenceExplainer
}

func NewExplainSentence(courseRepo repo.CourseRepo, explainer port.SentenceExplainer) *ExplainSentence {
	return &ExplainSentence{courseRepo: courseRepo, explainer: explainer}
}

func (uc *ExplainSentence) Exec(ctx context.Context, userID, courseID, rawSentence, sentenceTranslation string) (string, error) {
	course, err := uc.courseRepo.FindByID(ctx, userID, courseID)
	if err != nil {
		return "", err
	}

	sentence, err := valueobject.NewSentence(rawSentence)
	if err != nil {
		return "", err
	}
	if sentence.Text != rawSentence {
		return "", errs.WithHint(errs.ErrInvalidInput, "sentence is not normalized, send it as returned by /analyze")
	}

	explanation, err := uc.explainer.Explain(ctx, *sentence, strings.TrimSpace(sentenceTranslation), course.LearningLanguage, course.BaseLanguage)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", fmt.Errorf("%w: %v", errs.ErrUpstream, err)
	}
	return explanation, nil
}
