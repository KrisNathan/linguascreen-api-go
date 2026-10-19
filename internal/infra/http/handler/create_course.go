package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type createCourseRequest struct {
	BaseLanguage     string `json:"base_language" binding:"required"`
	LearningLanguage string `json:"learning_language" binding:"required"`
}

type CreateCourse struct {
	uc *usecase.CreateCourse
}

func NewCreateCourse(uc *usecase.CreateCourse) *CreateCourse {
	return &CreateCourse{uc: uc}
}

func (h *CreateCourse) Handle(c *gin.Context) {
	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	course, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), req.BaseLanguage, req.LearningLanguage)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewCourse(course))
}
