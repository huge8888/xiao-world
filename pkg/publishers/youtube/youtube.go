package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
)

// Publisher handles YouTube publishing
type Publisher struct {
	config     *configs.YouTubeConfig
	httpClient *http.Client
	enabled    bool
}

// NewPublisher creates a new YouTube publisher
func NewPublisher(cfg *configs.YouTubeConfig) *Publisher {
	return &Publisher{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 300 * time.Second, // 5 minutes for large video uploads
		},
		enabled: cfg != nil && cfg.Enabled && cfg.AccessToken != "",
	}
}

// GetName returns the publisher name
func (p *Publisher) GetName() string {
	return "YouTube"
}

// IsEnabled returns whether the publisher is enabled
func (p *Publisher) IsEnabled() bool {
	return p.enabled
}

// Publish publishes content to YouTube
func (p *Publisher) Publish(content *types.ProcessedContent) (*types.PublishResult, error) {
	result := &types.PublishResult{
		Platform:  types.PlatformYouTube,
		Timestamp: time.Now(),
	}

	if !p.enabled {
		result.Success = false
		result.Error = "YouTube publisher is not enabled or configured"
		return result, fmt.Errorf("publisher not enabled")
	}

	// YouTube only supports video content
	switch content.Type {
	case types.ContentTypeVideo:
		return p.publishVideo(content, result)
	case types.ContentTypeImage:
		result.Success = false
		result.Error = "YouTube only supports video content. Images not supported."
		return result, fmt.Errorf("images not supported")
	case types.ContentTypeText:
		result.Success = false
		result.Error = "YouTube requires video content. Text-only posts not supported."
		return result, fmt.Errorf("text-only not supported")
	default:
		result.Success = false
		result.Error = "unsupported content type"
		return result, fmt.Errorf("unsupported content type: %s", content.Type)
	}
}

// publishVideo publishes video to YouTube
func (p *Publisher) publishVideo(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	if len(content.MediaURLs) == 0 {
		result.Success = false
		result.Error = "no video URL provided"
		return result, fmt.Errorf("no video URL")
	}

	videoURL := content.MediaURLs[0]

	// Step 1: Download video
	videoData, err := p.downloadVideo(videoURL)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to download video: %v", err)
		return result, err
	}

	// Step 2: Prepare video metadata
	title := content.Title
	if title == "" {
		title = "Video from Xiaohongshu"
	}
	if len(title) > 100 {
		title = title[:100] // YouTube limit
	}

	description := content.Description
	if len(description) > 5000 {
		description = description[:5000] // YouTube limit
	}

	// Step 3: Upload video to YouTube
	videoID, videoURL, err := p.uploadVideo(videoData, title, description, content.Tags)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to upload video: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = videoID
	result.PostURL = videoURL

	return result, nil
}

// downloadVideo downloads video from URL
func (p *Publisher) downloadVideo(url string) ([]byte, error) {
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
		return nil, fmt.Errorf("failed to read video data: %w", err)
	}

	return data, nil
}

// uploadVideo uploads video to YouTube using YouTube Data API v3
func (p *Publisher) uploadVideo(videoData []byte, title, description string, tags []string) (string, string, error) {
	// Refresh access token if needed
	if err := p.refreshAccessToken(); err != nil {
		return "", "", fmt.Errorf("failed to refresh token: %w", err)
	}

	// YouTube Data API v3 endpoint
	apiURL := "https://www.googleapis.com/upload/youtube/v3/videos?uploadType=multipart&part=snippet,status"

	// Create video metadata
	metadata := map[string]interface{}{
		"snippet": map[string]interface{}{
			"title":       title,
			"description": description,
			"tags":        tags,
			"categoryId":  "22", // People & Blogs category
		},
		"status": map[string]interface{}{
			"privacyStatus": "public", // public, private, or unlisted
		},
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add metadata part
	metadataPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"application/json; charset=UTF-8"},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create metadata part: %w", err)
	}
	metadataPart.Write(metadataJSON)

	// Add video part
	videoPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"video/*"},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create video part: %w", err)
	}
	videoPart.Write(videoData)

	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.AccessToken))

	// Upload video
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API error: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	var uploadResp struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	videoID := uploadResp.ID
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	return videoID, videoURL, nil
}

// refreshAccessToken refreshes the OAuth 2.0 access token using refresh token
func (p *Publisher) refreshAccessToken() error {
	if p.config.RefreshToken == "" {
		// No refresh token, assume access token is still valid
		return nil
	}

	tokenURL := "https://oauth2.googleapis.com/token"

	reqBody := map[string]string{
		"client_id":     p.config.ClientID,
		"client_secret": p.config.ClientSecret,
		"refresh_token": p.config.RefreshToken,
		"grant_type":    "refresh_token",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh failed: %s (status: %d)", string(body), resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token,omitempty"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Update access token
	p.config.AccessToken = tokenResp.AccessToken

	// Update refresh token if a new one was provided
	if tokenResp.RefreshToken != "" {
		p.config.RefreshToken = tokenResp.RefreshToken
	}

	return nil
}
