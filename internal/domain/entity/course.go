package entity

import (
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

type Course struct {
	ID               string
	UserID           string
	BaseLanguage     valueobject.Language
	LearningLanguage valueobject.Language
	CreatedAt        time.Time
}
