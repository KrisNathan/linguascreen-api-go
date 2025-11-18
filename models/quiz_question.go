package models

type QuizQuestion struct {
	MemoryID int
	Question string
	Choices  []string
	AnswerID int
}
