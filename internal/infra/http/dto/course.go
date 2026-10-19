package dto

import (
	"time"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/model"
)

type Course struct {
	CourseID         string `json:"course_id"`
	BaseLanguage     string `json:"base_language"`
	LearningLanguage string `json:"learning_language"`
	CreatedAt        string `json:"created_at"`
}

func NewCourse(c *model.Course) Course {
	return Course{
		CourseID:         c.ID,
		BaseLanguage:     c.BaseLanguage.String(),
		LearningLanguage: c.LearningLanguage.String(),
		CreatedAt:        c.CreatedAt.UTC().Format(time.RFC3339),
	}
}
