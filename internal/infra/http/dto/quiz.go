package dto

import (
	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/aggregate"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type QuizScore struct {
	Correct  int `json:"correct"`
	Answered int `json:"answered"`
	Total    int `json:"total"`
}

type QuizQuestion struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Sentence string   `json:"sentence,omitempty"`
	Span     *[2]int  `json:"span,omitempty"`
	DictForm string   `json:"dict_form,omitempty"`
	Reading  string   `json:"reading,omitempty"`
	Options  []string `json:"options"`
}

type QuizSession struct {
	SessionID *string        `json:"session_id"`
	Score     QuizScore      `json:"score"`
	Questions []QuizQuestion `json:"questions"`
}

type QuizVocab struct {
	DictForm    string `json:"dict_form"`
	Reading     string `json:"reading"`
	Translation string `json:"translation"`
}

type QuizAnswer struct {
	Correct             bool      `json:"correct"`
	CorrectIndex        int       `json:"correct_index"`
	Vocab               QuizVocab `json:"vocab"`
	SentenceTranslation *string   `json:"sentence_translation,omitempty"`
	Score               QuizScore `json:"score"`
	SessionCompleted    bool      `json:"session_completed"`
}

func NewQuizSession(r *usecase.StartQuizResult) QuizSession {
	s := QuizSession{Score: newQuizScore(r.Score), Questions: make([]QuizQuestion, len(r.Questions))}
	if r.SessionID != "" {
		s.SessionID = &r.SessionID
	}
	for i, q := range r.Questions {
		s.Questions[i] = newQuizQuestion(q)
	}
	return s
}

func newQuizQuestion(q *aggregate.QuestionDetail) QuizQuestion {
	d := QuizQuestion{ID: q.ID, Type: q.Type.String(), Options: q.Options}
	switch q.Type {
	case valueobject.QuestionTypeFillBlank:
		d.Sentence = q.Occurrence.Sentence
		d.Span = &[2]int{q.Occurrence.Span.Start, q.Occurrence.Span.End}
	case valueobject.QuestionTypeTranslation:
		d.DictForm, d.Reading = q.Vocab.DictForm, q.Vocab.Reading
	case valueobject.QuestionTypeReading:
		d.DictForm = q.Vocab.DictForm
	}
	return d
}

func NewQuizAnswer(r *usecase.SubmitQuizResult) QuizAnswer {
	q := r.Question
	a := QuizAnswer{
		Correct:          *q.Correct,
		CorrectIndex:     q.CorrectIdx,
		Vocab:            QuizVocab{DictForm: q.Vocab.DictForm, Reading: q.Vocab.Reading, Translation: q.Vocab.Translation},
		Score:            newQuizScore(r.Score),
		SessionCompleted: r.SessionCompleted,
	}
	if q.Occurrence != nil {
		a.SentenceTranslation = &q.Occurrence.SentenceTranslation
	}
	return a
}

func newQuizScore(s valueobject.QuizScore) QuizScore {
	return QuizScore{Correct: s.Correct, Answered: s.Answered, Total: s.Total}
}
