package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type listCourseResponse struct {
	Courses []dto.Course `json:"courses"`
}

type ListCourse struct {
	uc *usecase.ListCourse
}

func NewListCourse(uc *usecase.ListCourse) *ListCourse {
	return &ListCourse{uc: uc}
}

func (h *ListCourse) Handle(c *gin.Context) {
	courses, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		c.Error(err)
		return
	}

	res := make([]dto.Course, len(courses))
	for i, co := range courses {
		res[i] = dto.NewCourse(co)
	}

	c.JSON(http.StatusOK, listCourseResponse{Courses: res})
}
