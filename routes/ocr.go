package routes

import (
	"net/http"
	"reflect"

	"linguascreen/models"
	"linguascreen/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("lang", validateLang)
}

func validateLang(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.String {
		return true
	}
	if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Interface {
		for i := 0; i < field.Len(); i++ {
			if field.Index(i).Elem().Kind() != reflect.String {
				return false
			}
		}
		return true
	}
	return false
}

type OcrRoutes struct {
	ocrService *services.OCRService
}

func NewOcrRoutes() *OcrRoutes {
	return &OcrRoutes{
		ocrService: services.NewOCRService(),
	}
}

type PostOcrRequest struct {
	Base64Image         string      `json:"base64Image" binding:"required"`
	Lang                interface{} `json:"lang" binding:"required" validate:"lang"`
	ConfidenceThreshold float64     `json:"confidenceThreshold,omitempty"`
}

type PostOcrResponse struct {
	Text string       `json:"text"`
	Page *models.Page `json:"page"`
}

func (r *OcrRoutes) Post(c *gin.Context) {
	var req PostOcrRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var langs []string
	switch v := req.Lang.(type) {
	case string:
		langs = []string{v}
	case []interface{}:
		for _, l := range v {
			langs = append(langs, l.(string))
		}
	}

	page, err := r.ocrService.ExtractTextFromBase64(req.Base64Image, langs, req.ConfidenceThreshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := PostOcrResponse{
		Text: page.Text,
		Page: page,
	}

	c.IndentedJSON(http.StatusOK, response)
}
