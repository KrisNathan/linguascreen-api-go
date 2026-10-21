package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type submitQuizRequest struct {
	QuestionID string `json:"question_id" binding:"required"`
	// Pointer because required rejects 0, which is a valid index.
	OptionIndex *int `json:"option_index" binding:"required"`
}

type SubmitQuiz struct {
	uc *usecase.SubmitQuiz
}

func NewSubmitQuiz(uc *usecase.SubmitQuiz) *SubmitQuiz {
	return &SubmitQuiz{uc: uc}
}

func (h *SubmitQuiz) Handle(c *gin.Context) {
	var req submitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	res, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), req.QuestionID, *req.OptionIndex)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewQuizAnswer(res))
}
