package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type startQuizRequest struct {
	CourseID string `json:"course_id" binding:"required"`
}

type StartQuiz struct {
	uc *usecase.StartQuiz
}

func NewStartQuiz(uc *usecase.StartQuiz) *StartQuiz {
	return &StartQuiz{uc: uc}
}

func (h *StartQuiz) Handle(c *gin.Context) {
	var req startQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	res, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), req.CourseID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewQuizSession(res))
}
