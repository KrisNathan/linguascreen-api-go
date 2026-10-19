package usecase

import (
	"context"
	"fmt"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/port"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

type AnalyzeSentence struct {
	courseRepo repo.CourseRepo
	analyzer   port.SentenceAnalyzer
}

func NewAnalyzeSentence(courseRepo repo.CourseRepo, analyzer port.SentenceAnalyzer) *AnalyzeSentence {
	return &AnalyzeSentence{courseRepo: courseRepo, analyzer: analyzer}
}

func (uc *AnalyzeSentence) Exec(ctx context.Context, userID, courseID, rawSentence string) (*valueobject.Sentence, *aggregate.Analysis, error) {
	sentence, err := valueobject.NewSentence(rawSentence)
	if err != nil {
		return nil, nil, err
	}

	course, err := uc.courseRepo.FindByID(ctx, userID, courseID)
	if err != nil {
		return nil, nil, err
	}

	for attempt := 1; ; attempt++ {
		analysis, err := uc.analyzer.Analyze(ctx, *sentence, course.LearningLanguage, course.BaseLanguage)
		if err == nil {
			err = analysis.Align(*sentence)
		}
		if err == nil {
			return sentence, analysis, nil
		}
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		if attempt == 2 {
			return nil, nil, fmt.Errorf("%w: %v", errs.ErrUpstream, err)
		}
	}
}
