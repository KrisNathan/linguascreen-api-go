package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type GetVocab struct {
	uc *usecase.GetVocab
}

func NewGetVocab(uc *usecase.GetVocab) *GetVocab {
	return &GetVocab{uc: uc}
}

func (h *GetVocab) Handle(c *gin.Context) {
	v, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), c.Param("vocab_id"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewVocabDetail(v))
}
