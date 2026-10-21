package aggregate

import "github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"

type Token struct {
	Surface                        string
	Span                           valueobject.Span // zero until aligned
	Saveable                       bool
	DictForm, Reading, Translation string
	Parts                          []Part
}
