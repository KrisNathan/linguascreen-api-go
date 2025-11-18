package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"linguascreen/models"
)

type OCRService struct {
	// Add fields if necessary
}

type OCROutput struct {
	Text string
}

func NewOCRService() *OCRService {
	return &OCRService{}
}

type OCRRequest struct {
	Base64Image string `json:"base64Image"`
	Lang        string `json:"lang"`
}

func (s *OCRService) ExtractTextFromBase64(base64Image, lang string) (*models.Page, error) {
	// Prepare the request payload
	reqPayload := OCRRequest{
		Base64Image: base64Image,
		Lang:        lang,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP POST request to the Node.js service
	resp, err := http.Post("http://localhost:3000/upload", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call OCR service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OCR service returned status %d", resp.StatusCode)
	}

	// Decode the response
	var page models.Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &page, nil
}
