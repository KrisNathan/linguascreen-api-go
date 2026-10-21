package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type signOutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SignOut struct {
	uc *usecase.SignOut
}

func NewSignOut(uc *usecase.SignOut) *SignOut {
	return &SignOut{uc: uc}
}

func (h *SignOut) Handle(c *gin.Context) {
	var req signOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.uc.Exec(c.Request.Context(), req.RefreshToken); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
