package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type DeleteVocab struct {
	uc *usecase.DeleteVocab
}

func NewDeleteVocab(uc *usecase.DeleteVocab) *DeleteVocab {
	return &DeleteVocab{uc: uc}
}

func (h *DeleteVocab) Handle(c *gin.Context) {
	if err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), c.Param("vocab_id")); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
