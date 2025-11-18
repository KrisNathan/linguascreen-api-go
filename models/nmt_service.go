package models

type NMTService interface {
	TranslateText(text string, targetLang string) (string, error)
}
