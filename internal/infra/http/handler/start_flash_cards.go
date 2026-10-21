package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type startFlashCardsRequest struct {
	CourseID string `json:"course_id" binding:"required"`
}

type StartFlashCards struct {
	uc *usecase.StartFlashCards
}

func NewStartFlashCards(uc *usecase.StartFlashCards) *StartFlashCards {
	return &StartFlashCards{uc: uc}
}

func (h *StartFlashCards) Handle(c *gin.Context) {
	var req startFlashCardsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	vocabs, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), req.CourseID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewFlashCardDeck(vocabs))
}
