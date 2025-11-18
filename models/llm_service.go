package models

type LLMService interface {
	Explain(sentence string, targetLang string, userLang string, nmtTranslation string, memories []Memory) (*ExplainResult, error)
	RearrangeText(text string) (string, error)
}
