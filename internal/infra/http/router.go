package http

import (
	"github.com/KrisNathan/linguascreen-api-go/internal/application/port"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/handler"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	SignIn          *handler.SignIn
	RefreshToken    *handler.RefreshToken
	SignOut         *handler.SignOut
	CreateCourse    *handler.CreateCourse
	ListCourse      *handler.ListCourse
	AnalyzeSentence *handler.AnalyzeSentence
	ExplainSentence *handler.ExplainSentence
	SaveVocab       *handler.SaveVocab
	ListVocabs      *handler.ListVocabs
	GetVocab        *handler.GetVocab
	UpdateVocab     *handler.UpdateVocab
	DeleteVocab     *handler.DeleteVocab
	StartQuiz       *handler.StartQuiz
	SubmitQuiz      *handler.SubmitQuiz
	StartFlashCards *handler.StartFlashCards
	SubmitFlashCard *handler.SubmitFlashCard
}

func NewRouter(h Handlers, tokenService port.TokenService, isProd bool) *gin.Engine {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.MapErrors())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth := router.Group("/auth")
	auth.POST("/google", h.SignIn.Handle)
	auth.POST("/refresh", h.RefreshToken.Handle)
	auth.POST("/sign-out", h.SignOut.Handle)

	api := router.Group("/", middleware.RequiresAuth(tokenService))
	api.POST("/courses", h.CreateCourse.Handle)
	api.GET("/courses", h.ListCourse.Handle)
	api.POST("/analyze", h.AnalyzeSentence.Handle)
	api.POST("/explain", h.ExplainSentence.Handle)
	api.POST("/vocabs", h.SaveVocab.Handle)
	api.GET("/courses/:course_id/vocabs", h.ListVocabs.Handle)
	api.GET("/vocabs/:vocab_id", h.GetVocab.Handle)
	api.PATCH("/vocabs/:vocab_id", h.UpdateVocab.Handle)
	api.DELETE("/vocabs/:vocab_id", h.DeleteVocab.Handle)
	api.POST("/quiz/start", h.StartQuiz.Handle)
	api.POST("/quiz/submit", h.SubmitQuiz.Handle)
	api.POST("/flash-card/start", h.StartFlashCards.Handle)
	api.POST("/flash-card/submit", h.SubmitFlashCard.Handle)

	return router
}
