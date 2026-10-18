package routes

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type signOutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SignOutHandler struct {
	uc *usecase.SignOut
}

func NewSignOutHandler(uc *usecase.SignOut) *SignOutHandler {
	return &SignOutHandler{uc: uc}
}

func (h *SignOutHandler) SignOutRoute(c *gin.Context) {
	var req signOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.uc.Exec(c.Request.Context(), req.RefreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}
