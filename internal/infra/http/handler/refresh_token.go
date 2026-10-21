package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshTokenResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresIn int    `json:"access_token_expires_in"`
	RefreshToken         string `json:"refresh_token"`
}

type RefreshToken struct {
	uc *usecase.RefreshToken
}

func NewRefreshToken(uc *usecase.RefreshToken) *RefreshToken {
	return &RefreshToken{uc: uc}
}

func (h *RefreshToken) Handle(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	pair, err := h.uc.Exec(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, refreshTokenResponse{
		AccessToken:          pair.AccessToken,
		AccessTokenExpiresIn: int(pair.AccessExpiresIn.Seconds()),
		RefreshToken:         pair.RefreshToken,
	})
}
