package models

type OCRService interface {
	ExtractTextFromBase64(base64Image string, lang []string, confidenceThreshold float64) (*Page, error)
}
