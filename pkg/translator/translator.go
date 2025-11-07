package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Translator interface for translation services
type Translator interface {
	Translate(text, sourceLang, targetLang string) (string, error)
	TranslateBatch(texts []string, sourceLang, targetLang string) ([]string, error)
}

// GoogleTranslator uses Google Translate API
type GoogleTranslator struct {
	apiKey     string
	httpClient *http.Client
}

// NewGoogleTranslator creates a new Google Translator
func NewGoogleTranslator(apiKey string) *GoogleTranslator {
	return &GoogleTranslator{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Translate translates text using Google Translate
func (t *GoogleTranslator) Translate(text, sourceLang, targetLang string) (string, error) {
	// If no API key, use free service (limited)
	if t.apiKey == "" {
		return t.translateFree(text, sourceLang, targetLang)
	}

	apiURL := fmt.Sprintf("https://translation.googleapis.com/language/translate/v2?key=%s", t.apiKey)

	reqBody := map[string]interface{}{
		"q":      []string{text},
		"source": sourceLang,
		"target": targetLang,
		"format": "text",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Data.Translations) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	return result.Data.Translations[0].TranslatedText, nil
}

// TranslateBatch translates multiple texts
func (t *GoogleTranslator) TranslateBatch(texts []string, sourceLang, targetLang string) ([]string, error) {
	if t.apiKey == "" {
		var results []string
		for _, text := range texts {
			translated, err := t.translateFree(text, sourceLang, targetLang)
			if err != nil {
				return nil, err
			}
			results = append(results, translated)
		}
		return results, nil
	}

	apiURL := fmt.Sprintf("https://translation.googleapis.com/language/translate/v2?key=%s", t.apiKey)

	reqBody := map[string]interface{}{
		"q":      texts,
		"source": sourceLang,
		"target": targetLang,
		"format": "text",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var translations []string
	for _, t := range result.Data.Translations {
		translations = append(translations, t.TranslatedText)
	}

	return translations, nil
}

// translateFree uses free Google Translate (limited, for development)
func (t *GoogleTranslator) translateFree(text, sourceLang, targetLang string) (string, error) {
	baseURL := "https://translate.googleapis.com/translate_a/single"

	params := url.Values{}
	params.Add("client", "gtx")
	params.Add("sl", sourceLang)
	params.Add("tl", targetLang)
	params.Add("dt", "t")
	params.Add("q", text)

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := t.httpClient.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	translations, ok := result[0].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	var translatedParts []string
	for _, part := range translations {
		if partArray, ok := part.([]interface{}); ok && len(partArray) > 0 {
			if translatedText, ok := partArray[0].(string); ok {
				translatedParts = append(translatedParts, translatedText)
			}
		}
	}

	return strings.Join(translatedParts, ""), nil
}
