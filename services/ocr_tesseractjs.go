package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	Base64Image string   `json:"base64Image"`
	Lang        []string `json:"lang"`
}

// filter words by confidence threshold (also reconstructs texts)
func filterWordsByConfidence(page *models.Page, threshold float64) {
	for i := range page.Blocks {
		block := &page.Blocks[i]
		blockTotalConfidence := 0.0

		for j := range block.Paragraphs {
			para := &block.Paragraphs[j]
			paraTotalConfidence := 0.0

			for k := range para.Lines {
				line := &para.Lines[k]
				var filteredWords []models.Word
				lineTotalConfidence := 0.0

				for _, word := range line.Words {
					if word.Confidence >= threshold {
						filteredWords = append(filteredWords, word)
						lineTotalConfidence += word.Confidence
					}
				}

				line.Words = filteredWords
				// avg confidence
				if len(line.Words) > 0 {
					line.Confidence = lineTotalConfidence / float64(len(line.Words))
				} else {
					line.Confidence = 0.0
				}
				paraTotalConfidence += line.Confidence

				// rebuild line.Text
				var texts []string
				for _, w := range line.Words {
					texts = append(texts, w.Text)
				}
				line.Text = strings.Join(texts, " ")
			}

			// avg confidence
			if len(para.Lines) > 0 {
				para.Confidence = paraTotalConfidence / float64(len(para.Lines))
			} else {
				para.Confidence = 0.0
			}
			blockTotalConfidence += para.Confidence

			// rebuild para.Text
			var lineTexts []string
			for _, l := range para.Lines {
				lineTexts = append(lineTexts, l.Text)
			}
			para.Text = strings.Join(lineTexts, "\n")
		}

		// avg confidence
		if len(block.Paragraphs) > 0 {
			block.Confidence = blockTotalConfidence / float64(len(block.Paragraphs))
		} else {
			block.Confidence = 0.0
		}

		// rebuild block.Text
		var paraTexts []string
		for _, p := range block.Paragraphs {
			paraTexts = append(paraTexts, p.Text)
		}
		block.Text = strings.Join(paraTexts, "\n\n")
	}

	// rebuild page.Text
	var blockTexts []string
	for _, b := range page.Blocks {
		blockTexts = append(blockTexts, b.Text)
	}
	page.Text = strings.Join(blockTexts, "\n\n")
}

func (s *OCRService) ExtractTextFromBase64(base64Image string, lang []string, confidenceThreshold float64) (*models.Page, error) {
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

	filterWordsByConfidence(&page, confidenceThreshold)

	return &page, nil
}
