package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AITranslator ใช้ AI Services (ChatGPT, Claude, Gemini) แปลภาษา
type AITranslator struct {
	provider   string // "openai", "anthropic", "google"
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewAITranslator สร้าง AI Translator ใหม่
// provider: "openai" (ChatGPT), "anthropic" (Claude), "google" (Gemini)
func NewAITranslator(provider, apiKey, model string) *AITranslator {
	if model == "" {
		// ใช้ model เริ่มต้นตามแต่ละ provider
		switch provider {
		case "openai":
			model = "gpt-4o-mini" // ถูกและเร็ว
		case "anthropic":
			model = "claude-3-haiku-20240307" // ถูกและเร็ว
		case "google":
			model = "gemini-1.5-flash" // ถูกและเร็ว
		}
	}

	return &AITranslator{
		provider: provider,
		apiKey:   apiKey,
		model:    model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Translate แปลข้อความด้วย AI
func (t *AITranslator) Translate(text, sourceLang, targetLang string) (string, error) {
	switch t.provider {
	case "openai":
		return t.translateWithOpenAI(text, sourceLang, targetLang)
	case "anthropic":
		return t.translateWithClaude(text, sourceLang, targetLang)
	case "google":
		return t.translateWithGemini(text, sourceLang, targetLang)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", t.provider)
	}
}

// TranslateBatch แปลหลายข้อความพร้อมกัน
func (t *AITranslator) TranslateBatch(texts []string, sourceLang, targetLang string) ([]string, error) {
	var results []string
	for _, text := range texts {
		translated, err := t.Translate(text, sourceLang, targetLang)
		if err != nil {
			return nil, err
		}
		results = append(results, translated)
	}
	return results, nil
}

// translateWithOpenAI แปลด้วย ChatGPT
func (t *AITranslator) translateWithOpenAI(text, sourceLang, targetLang string) (string, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"

	prompt := fmt.Sprintf("Translate the following %s text to %s. Return ONLY the translated text, no explanations:\n\n%s", sourceLang, targetLang, text)

	reqBody := map[string]interface{}{
		"model": t.model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3,
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	return result.Choices[0].Message.Content, nil
}

// translateWithClaude แปลด้วย Claude
func (t *AITranslator) translateWithClaude(text, sourceLang, targetLang string) (string, error) {
	apiURL := "https://api.anthropic.com/v1/messages"

	prompt := fmt.Sprintf("Translate the following %s text to %s. Return ONLY the translated text, no explanations:\n\n%s", sourceLang, targetLang, text)

	reqBody := map[string]interface{}{
		"model": t.model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 1024,
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
	req.Header.Set("x-api-key", t.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	return result.Content[0].Text, nil
}

// translateWithGemini แปลด้วย Gemini
func (t *AITranslator) translateWithGemini(text, sourceLang, targetLang string) (string, error) {
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", t.model, t.apiKey)

	prompt := fmt.Sprintf("Translate the following %s text to %s. Return ONLY the translated text, no explanations:\n\n%s", sourceLang, targetLang, text)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
