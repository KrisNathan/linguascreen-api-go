package middleware

import (
	"errors"
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"github.com/gin-gonic/gin"
)

// publicErrors maps errors that are safe to show clients to a status code.
// everything else becomes 500
var publicErrors = []struct {
	err    error
	status int
}{
	{errs.ErrInvalidInput, http.StatusBadRequest},
	{errs.ErrUnauthorized, http.StatusUnauthorized},
	{errs.ErrNotFound, http.StatusNotFound},
}

func MapErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		last := c.Errors.Last()
		// Handlers mark ShouldBind* failures with gin.ErrorTypeBind. The bind error itself only goes to the logs.
		if last.IsType(gin.ErrorTypeBind) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errs.ErrInvalidInput.Error(),
				"hint":  "request body must be JSON with all required fields",
			})
			return
		}

		err := last.Err
		for _, pe := range publicErrors {
			if errors.Is(err, pe.err) {
				body := gin.H{"error": pe.err.Error()}
				if hint := errs.Hint(err); hint != "" {
					body["hint"] = hint
				}
				c.JSON(pe.status, body)
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
