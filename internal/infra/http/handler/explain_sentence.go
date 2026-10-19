package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type explainSentenceRequest struct {
	CourseID            string `json:"course_id" binding:"required"`
	Sentence            string `json:"sentence" binding:"required"`
	SentenceTranslation string `json:"sentence_translation" binding:"required"`
}

type explainSentenceResponse struct {
	Explanation string `json:"explanation"`
}

type ExplainSentence struct {
	uc *usecase.ExplainSentence
}

func NewExplainSentence(uc *usecase.ExplainSentence) *ExplainSentence {
	return &ExplainSentence{uc: uc}
}

func (h *ExplainSentence) Handle(c *gin.Context) {
	var req explainSentenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	explanation, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), req.CourseID, req.Sentence, req.SentenceTranslation)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, explainSentenceResponse{Explanation: explanation})
}
