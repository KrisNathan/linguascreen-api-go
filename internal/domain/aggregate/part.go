package aggregate

import "github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"

type Part struct {
	Surface                        string
	Span                           valueobject.Span
	DictForm, Reading, Translation string
}
