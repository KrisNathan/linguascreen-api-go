package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type signInRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type signInResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresIn int          `json:"access_token_expires_in"`
	RefreshToken         string       `json:"refresh_token"`
	User                 userResponse `json:"user"`
}

type SignIn struct {
	uc *usecase.SignIn
}

func NewSignIn(uc *usecase.SignIn) *SignIn {
	return &SignIn{uc: uc}
}

func (h *SignIn) Handle(c *gin.Context) {
	var req signInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	res, err := h.uc.Exec(c.Request.Context(), req.IDToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		AccessToken:          res.Tokens.AccessToken,
		AccessTokenExpiresIn: int(res.Tokens.AccessExpiresIn.Seconds()),
		RefreshToken:         res.Tokens.RefreshToken,
		User:                 userResponse{ID: res.User.ID, Email: res.User.Email, Name: res.User.Name},
	})
}
