package routes

import (
	"errors"
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	errs "github.com/KrisNathan/linguascreen-api-go/internal/errors"
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

type SignInHandler struct {
	uc *usecase.SignIn
}

func NewSignInHandler(uc *usecase.SignIn) *SignInHandler {
	return &SignInHandler{uc: uc}
}

func (h *SignInHandler) SignInRoute(c *gin.Context) {
	var req signInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	res, err := h.uc.Exec(c.Request.Context(), req.IDToken)
	if errors.Is(err, errs.ErrUnauthorized) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		AccessToken:          res.Tokens.AccessToken,
		AccessTokenExpiresIn: int(res.Tokens.AccessExpiresIn.Seconds()),
		RefreshToken:         res.Tokens.RefreshToken,
		User:                 userResponse{ID: res.User.ID, Email: res.User.Email, Name: res.User.Name},
	})
}
