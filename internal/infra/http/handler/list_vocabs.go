package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type listVocabsQuery struct {
	Limit  int    `form:"limit"`
	Cursor string `form:"cursor"`
}

type ListVocabs struct {
	uc *usecase.ListVocabs
}

func NewListVocabs(uc *usecase.ListVocabs) *ListVocabs {
	return &ListVocabs{uc: uc}
}

func (h *ListVocabs) Handle(c *gin.Context) {
	var q listVocabsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	vocabs, next, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), c.Param("course_id"), q.Limit, q.Cursor)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewVocabList(vocabs, next))
}
