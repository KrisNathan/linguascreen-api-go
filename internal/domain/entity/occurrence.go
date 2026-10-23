package entity

import (
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type Occurrence struct {
	ID                  string
	Sentence            string
	SentenceTranslation string
	Span                valueobject.Span
	Surface             string
	Translation         string
	FromPart            bool
	CreatedAt           time.Time
}
