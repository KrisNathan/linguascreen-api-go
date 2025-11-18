package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// POST /v2/translate HTTP/2
// Host: api-free.deepl.com
// Authorization: DeepL-Auth-Key [yourAuthKey]
// User-Agent: YourApp/1.2.3
// Content-Length: 45
// Content-Type: application/json

// {"text":["Hello, world!"],"target_lang":"DE"}

// Response:
// {
//   "translations": [
//     {
//       "detected_source_language": "EN",
//       "text": "Hallo, Welt!"
//     }
//   ]
// }

type DeepLTranslateService struct {
	apiKey string
}

func NewDeepLTranslateService() *DeepLTranslateService {
	return &DeepLTranslateService{
		apiKey: os.Getenv("DEEPL_API_KEY"),
	}
}

type DeepLTranslateResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

func (s *DeepLTranslateService) TranslateText(text string, targetLang string) (string, error) {
	payload := map[string]interface{}{
		"text":        []string{text},
		"target_lang": targetLang,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("DeepL API error: %s", string(body))
	}

	var result DeepLTranslateResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if len(result.Translations) == 0 {
		return "", fmt.Errorf("no translations in response")
	}

	return result.Translations[0].Text, nil
}
