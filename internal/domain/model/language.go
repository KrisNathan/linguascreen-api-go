package model

import (
	"fmt"

	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
)

// supportedLanguages maps ISO 639-1 codes to the names used in LLM prompts.
var supportedLanguages = map[string]string{
	"en": "English",
	"ja": "Japanese",
	"zh": "Chinese",
}

// Language is a supported ISO 639-1 code.
// The zero value is invalid; construct with NewLanguage.
type Language struct {
	code string
}

func NewLanguage(code string) (Language, error) {
	if _, ok := supportedLanguages[code]; !ok {
		return Language{}, errs.WithHint(errs.ErrInvalidInput, fmt.Sprintf("%q is not a supported language", code))
	}
	return Language{code: code}, nil
}

func (l Language) String() string { return l.code }

// Name returns the English name of the language, e.g. "Japanese".
func (l Language) Name() string { return supportedLanguages[l.code] }
