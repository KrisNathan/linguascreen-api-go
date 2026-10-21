package valueobject

import "math/rand"

type QuestionType int

const (
	QuestionTypeUnknown QuestionType = iota
	QuestionTypeFillBlank
	QuestionTypeTranslation
	QuestionTypeReading
)

func QuestionTypeFromString(str string) QuestionType {
	switch str {
	case "fill_blank":
		return QuestionTypeFillBlank
	case "translation":
		return QuestionTypeTranslation
	case "reading":
		return QuestionTypeReading
	default:
		return QuestionTypeUnknown
	}
}

func ShuffledQuestionTypes() []QuestionType {
	types := []QuestionType{QuestionTypeFillBlank, QuestionTypeTranslation, QuestionTypeReading}
	rand.Shuffle(len(types), func(i, j int) { types[i], types[j] = types[j], types[i] })
	return types
}

func (t QuestionType) String() string {
	switch t {
	case QuestionTypeFillBlank:
		return "fill_blank"
	case QuestionTypeTranslation:
		return "translation"
	case QuestionTypeReading:
		return "reading"
	default:
		return "unknown"
	}
}
