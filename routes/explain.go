package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExplainRoutes struct{}

func NewExplainRoutes() *ExplainRoutes {
	return &ExplainRoutes{}
}

type PostExplainResponse struct {
	Explanation string `json:"explanation"`
}

func (r *ExplainRoutes) Post(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, PostExplainResponse{Explanation: "Sample explanation text"})
}
