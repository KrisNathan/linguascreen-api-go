package routes

import (
	"net/http"

	"linguascreen/models"
	"linguascreen/services"

	"github.com/gin-gonic/gin"
)

type OcrRoutes struct {
	ocrService *services.OCRService
}

func NewOcrRoutes() *OcrRoutes {
	return &OcrRoutes{
		ocrService: services.NewOCRService(),
	}
}

type PostOcrRequest struct {
	Base64Image string `json:"base64Image" binding:"required"`
	Lang        string `json:"lang" binding:"required"`
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

	page, err := r.ocrService.ExtractTextFromBase64(req.Base64Image, req.Lang)
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
