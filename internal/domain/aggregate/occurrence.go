package aggregate

import "github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"

type Occurrence struct {
	ID          string
	Span        valueobject.Span
	Surface     string
	Translation string
	FromPart    bool
}
