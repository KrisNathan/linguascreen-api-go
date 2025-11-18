package routes

import (
	"linguascreen/models"
	"linguascreen/repository"
	"math/rand"

	"github.com/gin-gonic/gin"
)

type QuizRoutes struct {
	memoryRepo *repository.MemoryRepository
}

func NewQuizRoutes(memoryRepo *repository.MemoryRepository) *QuizRoutes {
	return &QuizRoutes{memoryRepo: memoryRepo}
}

type QuizRequest struct {
	UserID int `uri:"user_id" binding:"required"`
}

type GetQuizResponse struct {
	Quiz models.Quiz
}

func (r *QuizRoutes) Get(c *gin.Context) {
	var req QuizRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	memories, err := r.memoryRepo.GetWeakestMemories(req.UserID, 10)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get memories"})
		return
	}

	// Generate quiz questions from memories (dummy implementation)
	var questions []models.QuizQuestion
	for _, memory := range memories {
		wrongs, err := r.memoryRepo.GetRandomMemoriesByContentType(memory.ContentType, memory.ID, 3)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get wrong options"})
			return
		}

		choices := []string{memory.Translation}
		for _, w := range wrongs {
			choices = append(choices, w.Translation)
		}

		// Shuffle the choices
		rand.Shuffle(len(choices), func(i, j int) {
			choices[i], choices[j] = choices[j], choices[i]
		})

		// Find the index of the correct answer
		answerID := 0
		for i, c := range choices {
			if c == memory.Translation {
				answerID = i
				break
			}
		}

		question := models.QuizQuestion{
			MemoryID: memory.ID,
			Question: "What is the translation of '" + memory.Title + "'?",
			Choices:  choices,
			AnswerID: answerID,
		}
		questions = append(questions, question)
	}

	quiz := models.Quiz{
		QuestionCount: len(questions),
		Questions:     questions,
	}

	response := GetQuizResponse{
		Quiz: quiz,
	}

	c.JSON(200, response)
}
