package models

type ExplainResult struct {
	GrammarExplanations []MemoryInput `json:"grammar_explanations" jsonschema_description:"A list of new grammar explanations relevant to the input"`
	PhraseExplanations  []MemoryInput `json:"phrase_explanations" jsonschema_description:"A list of new phrase explanations relevant to the input"`
	WordExplanations    []MemoryInput `json:"word_explanations" jsonschema_description:"A list of new word explanations relevant to the input"`
}
