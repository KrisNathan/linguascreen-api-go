package routes

import (
	"errors"
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	errs "github.com/KrisNathan/linguascreen-api-go/internal/errors"
	"github.com/gin-gonic/gin"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresIn int    `json:"access_token_expires_in"`
	RefreshToken         string `json:"refresh_token"`
}

type RefreshTokenHandler struct {
	uc *usecase.RefreshToken
}

func NewRefreshTokenHandler(uc *usecase.RefreshToken) *RefreshTokenHandler {
	return &RefreshTokenHandler{uc: uc}
}

func (h *RefreshTokenHandler) RefreshTokenRoute(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	pair, err := h.uc.Exec(c.Request.Context(), req.RefreshToken)
	if errors.Is(err, errs.ErrUnauthorized) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, refreshResponse{
		AccessToken:          pair.AccessToken,
		AccessTokenExpiresIn: int(pair.AccessExpiresIn.Seconds()),
		RefreshToken:         pair.RefreshToken,
	})
}
