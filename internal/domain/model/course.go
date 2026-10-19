package model

import "time"

type Course struct {
	ID               string
	UserID           string
	BaseLanguage     Language
	LearningLanguage Language
	CreatedAt        time.Time
}
