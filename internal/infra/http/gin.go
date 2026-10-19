package http

import (
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/routes"
	"github.com/gin-gonic/gin"
)

func NewRouter(signIn *routes.SignInHandler, refresh *routes.RefreshTokenHandler, signOut *routes.SignOutHandler, isProd bool) *gin.Engine {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/auth/google", signIn.SignInRoute)
	router.POST("/auth/refresh", refresh.RefreshTokenRoute)
	router.POST("/auth/sign-out", signOut.SignOutRoute)

	return router
}
