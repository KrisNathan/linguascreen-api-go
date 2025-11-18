package models

type Memory struct {
	ID     int
	UserID int

	ContentType     string // e.g., 'word', 'grammar', 'phrase'
	Title           string // e.g., the word or phrase itself or grammar name
	Translation     string
	Explanation     string
	ExampleSentence string

	MemoryStrength int // e.g., 0 to 100 scale
	LastReviewedAt string
	CreatedAt      string
}

type MemoryInput struct {
	ContentType     string
	Title           string
	Translation     string
	Explanation     string
	ExampleSentence string
}
