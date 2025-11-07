package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
)

// Publisher handles Twitter/X publishing
type Publisher struct {
	config     *configs.TwitterConfig
	httpClient *http.Client
	enabled    bool
}

// NewPublisher creates a new Twitter publisher
func NewPublisher(cfg *configs.TwitterConfig) *Publisher {
	return &Publisher{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		enabled: cfg != nil && cfg.Enabled && cfg.BearerToken != "",
	}
}

// GetName returns the publisher name
func (p *Publisher) GetName() string {
	return "Twitter/X"
}

// IsEnabled returns whether the publisher is enabled
func (p *Publisher) IsEnabled() bool {
	return p.enabled
}

// Publish publishes content to Twitter
func (p *Publisher) Publish(content *types.ProcessedContent) (*types.PublishResult, error) {
	result := &types.PublishResult{
		Platform:  types.PlatformTwitter,
		Timestamp: time.Now(),
	}

	if !p.enabled {
		result.Success = false
		result.Error = "Twitter publisher is not enabled or configured"
		return result, fmt.Errorf("publisher not enabled")
	}

	// Handle different content types
	switch content.Type {
	case types.ContentTypeText:
		return p.publishText(content, result)
	case types.ContentTypeImage:
		return p.publishWithImages(content, result)
	case types.ContentTypeVideo:
		result.Success = false
		result.Error = "Video publishing to Twitter not yet implemented"
		return result, fmt.Errorf("video not supported yet")
	default:
		result.Success = false
		result.Error = "unsupported content type"
		return result, fmt.Errorf("unsupported content type: %s", content.Type)
	}
}

// publishText publishes text-only tweet
func (p *Publisher) publishText(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	// Using Twitter API v2
	apiURL := "https://api.twitter.com/2/tweets"

	reqBody := map[string]interface{}{
		"text": content.Description,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to marshal request: %v", err)
		return result, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.BearerToken))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to send request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		result.Success = false
		result.Error = fmt.Sprintf("API error: %s (status: %d)", string(body), resp.StatusCode)
		return result, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	// Parse response
	var tweetResp struct {
		Data struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &tweetResp); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to parse response: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = tweetResp.Data.ID
	result.PostURL = fmt.Sprintf("https://twitter.com/i/web/status/%s", tweetResp.Data.ID)

	return result, nil
}

// publishWithImages publishes tweet with images
func (p *Publisher) publishWithImages(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	if len(content.MediaURLs) == 0 {
		// No images, publish as text
		return p.publishText(content, result)
	}

	// Twitter supports max 4 images per tweet
	maxImages := 4
	if len(content.MediaURLs) > maxImages {
		content.MediaURLs = content.MediaURLs[:maxImages]
	}

	// Step 1: Download and upload all images to get media IDs
	var mediaIDs []string
	for _, imageURL := range content.MediaURLs {
		mediaID, err := p.uploadImage(imageURL)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("failed to upload image: %v", err)
			return result, err
		}
		mediaIDs = append(mediaIDs, mediaID)
	}

	// Step 2: Create tweet with media IDs
	apiURL := "https://api.twitter.com/2/tweets"

	reqBody := map[string]interface{}{
		"text": content.Description,
		"media": map[string]interface{}{
			"media_ids": mediaIDs,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to marshal request: %v", err)
		return result, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.BearerToken))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to send request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		result.Success = false
		result.Error = fmt.Sprintf("API error: %s (status: %d)", string(body), resp.StatusCode)
		return result, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var tweetResp struct {
		Data struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &tweetResp); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to parse response: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = tweetResp.Data.ID
	result.PostURL = fmt.Sprintf("https://twitter.com/i/web/status/%s", tweetResp.Data.ID)

	return result, nil
}

// uploadImage downloads and uploads image to Twitter, returns media ID
func (p *Publisher) uploadImage(imageURL string) (string, error) {
	// Step 1: Download image
	imageData, err := p.downloadImage(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}

	// Step 2: Upload to Twitter media endpoint
	// Using Twitter API v1.1 for media upload (v2 doesn't support media upload yet)
	uploadURL := "https://upload.twitter.com/1.1/media/upload.json"

	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.BearerToken))
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed: %s (status: %d)", string(body), resp.StatusCode)
	}

	var uploadResp struct {
		MediaID       int64  `json:"media_id"`
		MediaIDString string `json:"media_id_string"`
	}

	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to parse upload response: %w", err)
	}

	return uploadResp.MediaIDString, nil
}

// downloadImage downloads image from URL
func (p *Publisher) downloadImage(url string) ([]byte, error) {
	resp, err := p.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return data, nil
}
