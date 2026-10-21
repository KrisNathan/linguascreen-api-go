package handler

import (
	"net/http"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/usecase"
	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/dto"
	"github.com/KrisNathan/linguascreen-api-go/internal/infra/http/middleware"
	"github.com/gin-gonic/gin"
)

type saveVocabRequest struct {
	CourseID            string `json:"course_id" binding:"required"`
	Sentence            string `json:"sentence" binding:"required"`
	SentenceTranslation string `json:"sentence_translation" binding:"required"`
	Span                []int  `json:"span" binding:"required,len=2"`
	FromPart            bool   `json:"from_part"`
	DictForm            string `json:"dict_form" binding:"required"`
	Reading             string `json:"reading"`
	Translation         string `json:"translation" binding:"required"`
}

type SaveVocab struct {
	uc *usecase.SaveVocab
}

func NewSaveVocab(uc *usecase.SaveVocab) *SaveVocab {
	return &SaveVocab{uc: uc}
}

func (h *SaveVocab) Handle(c *gin.Context) {
	var req saveVocabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	res, err := h.uc.Exec(c.Request.Context(), middleware.UserID(c), usecase.SaveVocabInput{
		CourseID:            req.CourseID,
		Sentence:            req.Sentence,
		SentenceTranslation: req.SentenceTranslation,
		Span:                valueobject.Span{Start: req.Span[0], End: req.Span[1]},
		FromPart:            req.FromPart,
		DictForm:            req.DictForm,
		Reading:             req.Reading,
		Translation:         req.Translation,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewSavedVocab(res))
}
