package service

import (
	"context"
	"math/rand"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

const (
	QuizSize           = 10
	distractorPoolSize = 50
	wrongOptionCount   = 3
)

type QuizService struct {
	vocabRepo repo.VocabRepo
}

func NewQuizService(vocabRepo repo.VocabRepo) *QuizService {
	return &QuizService{vocabRepo: vocabRepo}
}

// NewSession doesn't save the session.
func (s *QuizService) NewSession(ctx context.Context, courseID string) (*aggregate.PracticeSession, error) {
	// Extra because vocabs with no eligible question type are skipped, see docs/use_cases/scheduling.md.
	vocabs, err := s.vocabRepo.ListForPractice(ctx, courseID, QuizSize*2)
	if err != nil {
		return nil, err
	}
	pool, err := s.vocabRepo.FindDistractorVocabs(ctx, courseID, distractorPoolSize)
	if err != nil {
		return nil, err
	}

	session := &aggregate.PracticeSession{CourseID: courseID, Questions: []*aggregate.QuestionDetail{}}
	for _, vocab := range vocabs {
		if len(session.Questions) == QuizSize {
			break
		}
		q, err := s.questionFor(ctx, vocab, pool)
		if err != nil {
			return nil, err
		}
		if q == nil {
			continue
		}
		q.Position = len(session.Questions)
		session.Questions = append(session.Questions, q)
	}
	return session, nil
}

func (s *QuizService) questionFor(ctx context.Context, vocab *aggregate.Vocab, pool []*aggregate.Vocab) (*aggregate.QuestionDetail, error) {
	for _, t := range valueobject.ShuffledQuestionTypes() {
		var q *aggregate.QuestionDetail
		switch t {
		case valueobject.QuestionTypeFillBlank:
			var err error
			if q, err = s.NewFillBlankQuestion(ctx, vocab, pool); err != nil {
				return nil, err
			}
		case valueobject.QuestionTypeTranslation:
			q = NewTranslationQuestion(vocab, pool)
		case valueobject.QuestionTypeReading:
			q = NewReadingQuestion(vocab, pool)
		}
		if q != nil {
			return q, nil
		}
	}
	return nil, nil
}

func (s *QuizService) NewFillBlankQuestion(ctx context.Context, vocab *aggregate.Vocab, pool []*aggregate.Vocab) (*aggregate.QuestionDetail, error) {
	wrong := pickWrongOptions(vocab, pool, func(v *aggregate.Vocab) string { return v.DictForm })
	if wrong == nil {
		return nil, nil
	}
	occ, err := s.vocabRepo.FindQuizOccurrence(ctx, vocab.ID)
	if err != nil || occ == nil {
		return nil, err
	}
	options, correctIdx := withCorrectOption(wrong, vocab.DictForm)
	return &aggregate.QuestionDetail{
		Question: &aggregate.Question{
			Type:         valueobject.QuestionTypeFillBlank,
			VocabID:      vocab.ID,
			OccurrenceID: &occ.ID,
			Options:      options,
			CorrectIdx:   correctIdx,
		},
		Vocab:      vocab,
		Occurrence: occ,
	}, nil
}

func NewTranslationQuestion(vocab *aggregate.Vocab, pool []*aggregate.Vocab) *aggregate.QuestionDetail {
	wrong := pickWrongOptions(vocab, pool, func(v *aggregate.Vocab) string { return v.Translation })
	if wrong == nil {
		return nil
	}
	options, correctIdx := withCorrectOption(wrong, vocab.Translation)
	return &aggregate.QuestionDetail{
		Question: &aggregate.Question{
			Type:       valueobject.QuestionTypeTranslation,
			VocabID:    vocab.ID,
			Options:    options,
			CorrectIdx: correctIdx,
		},
		Vocab: vocab,
	}
}

func NewReadingQuestion(vocab *aggregate.Vocab, pool []*aggregate.Vocab) *aggregate.QuestionDetail {
	// A kana-only word like よく would give the answer away.
	if vocab.Reading == "" || vocab.Reading == vocab.DictForm {
		return nil
	}
	wrong := pickWrongOptions(vocab, pool, func(v *aggregate.Vocab) string { return v.Reading })
	if wrong == nil {
		return nil
	}
	options, correctIdx := withCorrectOption(wrong, vocab.Reading)
	return &aggregate.QuestionDetail{
		Question: &aggregate.Question{
			Type:       valueobject.QuestionTypeReading,
			VocabID:    vocab.ID,
			Options:    options,
			CorrectIdx: correctIdx,
		},
		Vocab: vocab,
	}
}

func pickWrongOptions(vocab *aggregate.Vocab, pool []*aggregate.Vocab, field func(*aggregate.Vocab) string) []string {
	seen := map[string]bool{field(vocab): true}
	wrong := make([]string, 0, wrongOptionCount)
	for _, i := range rand.Perm(len(pool)) {
		v := pool[i]
		o := field(v)
		if v.ID == vocab.ID || o == "" || seen[o] {
			continue
		}
		seen[o] = true
		wrong = append(wrong, o)
		if len(wrong) == wrongOptionCount {
			return wrong
		}
	}
	return nil
}

func withCorrectOption(wrong []string, correct string) ([]string, int) {
	idx := rand.Intn(len(wrong) + 1)
	options := make([]string, 0, len(wrong)+1)
	options = append(options, wrong[:idx]...)
	options = append(options, correct)
	options = append(options, wrong[idx:]...)
	return options, idx
}
