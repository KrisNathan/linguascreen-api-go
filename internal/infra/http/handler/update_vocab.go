package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type updateVocabRequest struct {
	Translation string `json:"translation" binding:"required"`
}

type UpdateVocab struct {
	uc *usecase.UpdateVocab
}

func NewUpdateVocab(uc *usecase.UpdateVocab) *UpdateVocab {
	return &UpdateVocab{uc: uc}
}

func (h *UpdateVocab) Handle(c *gin.Context) {
	var req updateVocabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	v, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), c.Param("vocab_id"), req.Translation)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewVocabItem(v))
}
