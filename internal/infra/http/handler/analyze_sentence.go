package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type analyzeSentenceRequest struct {
	CourseID string `json:"course_id" binding:"required"`
	Sentence string `json:"sentence" binding:"required"`
}

type AnalyzeSentence struct {
	uc *usecase.AnalyzeSentence
}

func NewAnalyzeSentence(uc *usecase.AnalyzeSentence) *AnalyzeSentence {
	return &AnalyzeSentence{uc: uc}
}

func (h *AnalyzeSentence) Handle(c *gin.Context) {
	var req analyzeSentenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	sentence, analysis, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), req.CourseID, req.Sentence)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewAnalysis(sentence, analysis))
}
