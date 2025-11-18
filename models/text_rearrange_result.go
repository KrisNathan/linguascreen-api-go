package models

type TextRearrangeResult struct {
	ProcessedText string `json:"processed_text" jsonschema_description:"The text after preprocessing"`
}
