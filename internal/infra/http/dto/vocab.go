package dto

import (
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type SavedVocab struct {
	VocabID           string `json:"vocab_id"`
	OccurrenceID      string `json:"occurrence_id"`
	VocabCreated      bool   `json:"vocab_created"`
	OccurrenceCreated bool   `json:"occurrence_created"`
}

func NewSavedVocab(r *usecase.SaveVocabResult) SavedVocab {
	return SavedVocab{
		VocabID:           r.VocabID,
		OccurrenceID:      r.OccurrenceID,
		VocabCreated:      r.VocabCreated,
		OccurrenceCreated: r.OccurrenceCreated,
	}
}

type Progress struct {
	State  string  `json:"state"`
	Streak int     `json:"streak"`
	DueAt  *string `json:"due_at"`
}

type VocabItem struct {
	ID              string   `json:"id"`
	DictForm        string   `json:"dict_form"`
	Reading         string   `json:"reading"`
	Translation     string   `json:"translation"`
	OccurrenceCount int      `json:"occurrence_count"`
	LastSavedAt     string   `json:"last_saved_at"`
	Progress        Progress `json:"progress"`
}

type VocabList struct {
	Items      []VocabItem `json:"items"`
	NextCursor *string     `json:"next_cursor"`
}

type Occurrence struct {
	ID                  string `json:"id"`
	Sentence            string `json:"sentence"`
	SentenceTranslation string `json:"sentence_translation"`
	Surface             string `json:"surface"`
	Span                [2]int `json:"span"`
	Translation         string `json:"translation"`
	FromPart            bool   `json:"from_part"`
	CreatedAt           string `json:"created_at"`
}

type VocabDetail struct {
	ID          string       `json:"id"`
	DictForm    string       `json:"dict_form"`
	Reading     string       `json:"reading"`
	Translation string       `json:"translation"`
	LastSavedAt string       `json:"last_saved_at"`
	Progress    Progress     `json:"progress"`
	Occurrences []Occurrence `json:"occurrences"`
}

func NewVocabList(vocabs []*aggregate.Vocab, nextCursor string) VocabList {
	l := VocabList{Items: make([]VocabItem, len(vocabs))}
	for i, v := range vocabs {
		l.Items[i] = VocabItem{
			ID:              v.ID,
			DictForm:        v.DictForm,
			Reading:         v.Reading,
			Translation:     v.Translation,
			OccurrenceCount: v.OccurrenceCount,
			LastSavedAt:     formatTime(v.LastSavedAt),
			Progress:        newProgress(v.Progress),
		}
	}
	if nextCursor != "" {
		l.NextCursor = &nextCursor
	}
	return l
}

func NewVocabDetail(v *aggregate.Vocab) VocabDetail {
	d := VocabDetail{
		ID:          v.ID,
		DictForm:    v.DictForm,
		Reading:     v.Reading,
		Translation: v.Translation,
		LastSavedAt: formatTime(v.LastSavedAt),
		Progress:    newProgress(v.Progress),
		Occurrences: make([]Occurrence, len(v.Occurrences)),
	}
	for i, o := range v.Occurrences {
		d.Occurrences[i] = Occurrence{
			ID:                  o.ID,
			Sentence:            o.Sentence,
			SentenceTranslation: o.SentenceTranslation,
			Surface:             o.Surface,
			Span:                [2]int{o.Span.Start, o.Span.End},
			Translation:         o.Translation,
			FromPart:            o.FromPart,
			CreatedAt:           formatTime(o.CreatedAt),
		}
	}
	return d
}

func newProgress(p valueobject.Progress) Progress {
	out := Progress{State: p.State, Streak: p.Streak}
	if p.DueAt != nil {
		s := formatTime(*p.DueAt)
		out.DueAt = &s
	}
	return out
}

func formatTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }
