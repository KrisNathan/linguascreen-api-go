package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/repo"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

const (
	defaultVocabPageSize = 30
	maxVocabPageSize     = 100
)

type ListVocabs struct {
	courseRepo repo.CourseRepo
	vocabRepo  repo.VocabRepo
}

func NewListVocabs(courseRepo repo.CourseRepo, vocabRepo repo.VocabRepo) *ListVocabs {
	return &ListVocabs{courseRepo: courseRepo, vocabRepo: vocabRepo}
}

func (uc *ListVocabs) Exec(ctx context.Context, userID, courseID string, limit int, cursor string) ([]*aggregate.Vocab, string, error) {
	if _, err := uc.courseRepo.FindByID(ctx, userID, courseID); err != nil {
		return nil, "", err
	}

	if limit == 0 {
		limit = defaultVocabPageSize
	}
	if limit < 1 || limit > maxVocabPageSize {
		return nil, "", errs.WithHint(errs.ErrInvalidInput, fmt.Sprintf("limit must be between 1 and %d", maxVocabPageSize))
	}

	var after *repo.VocabCursor
	if cursor != "" {
		c, err := decodeVocabCursor(cursor)
		if err != nil {
			return nil, "", errs.WithHint(errs.ErrInvalidInput, "invalid cursor")
		}
		after = &c
	}

	vocabs, err := uc.vocabRepo.List(ctx, courseID, after, limit+1)
	if err != nil {
		return nil, "", err
	}

	next := ""
	if len(vocabs) > limit {
		vocabs = vocabs[:limit]
		last := vocabs[limit-1]
		next = encodeVocabCursor(repo.VocabCursor{LastSavedAt: last.LastSavedAt, ID: last.ID})
	}
	for _, v := range vocabs {
		v.Progress = valueobject.Progress{State: valueobject.ProgressNew} // TODO: derive from answered questions once Quiz exists
	}
	return vocabs, next, nil
}

// Postgres stores microseconds, so RFC3339Nano round-trips last_saved_at exactly.
func encodeVocabCursor(c repo.VocabCursor) string {
	return base64.RawURLEncoding.EncodeToString([]byte(c.LastSavedAt.UTC().Format(time.RFC3339Nano) + "|" + c.ID))
}

func decodeVocabCursor(s string) (repo.VocabCursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return repo.VocabCursor{}, err
	}
	at, id, ok := strings.Cut(string(b), "|")
	if !ok || id == "" {
		return repo.VocabCursor{}, fmt.Errorf("malformed cursor")
	}
	t, err := time.Parse(time.RFC3339Nano, at)
	if err != nil {
		return repo.VocabCursor{}, err
	}
	return repo.VocabCursor{LastSavedAt: t, ID: id}, nil
}
