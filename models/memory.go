package models

type Memory struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	ContentType     string `json:"content_type"` // e.g., 'word', 'grammar', 'phrase'
	Title           string `json:"title"`        // e.g., the word or phrase itself or grammar name
	Translation     string `json:"translation"`
	Explanation     string `json:"explanation"`
	ExampleSentence string `json:"example_sentence"`
	MemoryStrength  int    `json:"memory_strength"` // e.g., 0 to 100 scale
	LastReviewedAt  string `json:"last_reviewed_at"`
	CreatedAt       string `json:"created_at"`
}

type MemoryInput struct {
	ContentType     string `json:"content_type" jsonschema_description:"The type of content being explained, possible options: 'word', 'phrase', 'grammar'"`
	Title           string `json:"title" jsonschema_description:"The word, phrase, or grammar point being explained"`
	Translation     string `json:"translation" jsonschema_description:"The correct translation of the given sentence"`
	Explanation     string `json:"explanation" jsonschema_description:"A detailed explanation of the grammar point involved"`
	ExampleSentence string `json:"example_sentence" jsonschema_description:"An example sentence using the word, phrase, or grammar point"`
}
