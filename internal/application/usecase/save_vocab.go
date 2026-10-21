package usecase

import (
	"context"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

type SaveVocabInput struct {
	CourseID            string
	Sentence            string
	SentenceTranslation string
	Span                valueobject.Span
	FromPart            bool
	DictForm            string
	Reading             string
	Translation         string
}

type SaveVocabResult struct {
	VocabID           string
	OccurrenceID      string
	VocabCreated      bool
	OccurrenceCreated bool
}

type SaveVocab struct {
	courseRepo repo.CourseRepo
	vocabRepo  repo.VocabRepo
}

func NewSaveVocab(courseRepo repo.CourseRepo, vocabRepo repo.VocabRepo) *SaveVocab {
	return &SaveVocab{courseRepo: courseRepo, vocabRepo: vocabRepo}
}

func (uc *SaveVocab) Exec(ctx context.Context, userID string, in SaveVocabInput) (*SaveVocabResult, error) {
	course, err := uc.courseRepo.FindByID(ctx, userID, in.CourseID)
	if err != nil {
		return nil, err
	}

	sentence, err := valueobject.NewSentence(in.Sentence)
	if err != nil {
		return nil, err
	}
	// Spans only make sense against the exact text /analyze returned.
	if sentence.Text != in.Sentence {
		return nil, errs.WithHint(errs.ErrInvalidInput, "sentence is not normalized, send it as returned by /analyze")
	}

	surface, ok := in.Span.Slice(sentence.Text)
	if !ok {
		return nil, errs.WithHint(errs.ErrInvalidInput, "span is empty, out of range or splits a character")
	}

	vocab, err := aggregate.NewVocab(course.ID, course.LearningLanguage, in.DictForm, in.Reading, in.Translation)
	if err != nil {
		return nil, err
	}

	occ := aggregate.Occurrence{Span: in.Span, Surface: surface, Translation: vocab.Translation, FromPart: in.FromPart}
	vocabCreated, occCreated, err := uc.vocabRepo.Save(ctx, *sentence, in.SentenceTranslation, vocab, &occ)
	if err != nil {
		return nil, err
	}

	return &SaveVocabResult{
		VocabID:           vocab.ID,
		OccurrenceID:      occ.ID,
		VocabCreated:      vocabCreated,
		OccurrenceCreated: occCreated,
	}, nil
}
