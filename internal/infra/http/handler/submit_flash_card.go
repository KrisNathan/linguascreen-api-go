package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type submitFlashCardRequest struct {
	VocabID string `json:"vocab_id" binding:"required"`
	// Pointer because required rejects false.
	Correct *bool `json:"correct" binding:"required"`
}

type SubmitFlashCard struct {
	uc *usecase.SubmitFlashCard
}

func NewSubmitFlashCard(uc *usecase.SubmitFlashCard) *SubmitFlashCard {
	return &SubmitFlashCard{uc: uc}
}

func (h *SubmitFlashCard) Handle(c *gin.Context) {
	var req submitFlashCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), req.VocabID, *req.Correct); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
