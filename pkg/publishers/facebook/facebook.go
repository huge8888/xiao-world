package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
)

// Publisher handles Facebook publishing
type Publisher struct {
	config     *configs.FacebookConfig
	httpClient *http.Client
	enabled    bool
}

// NewPublisher creates a new Facebook publisher
func NewPublisher(cfg *configs.FacebookConfig) *Publisher {
	return &Publisher{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		enabled: cfg != nil && cfg.Enabled && cfg.AccessToken != "" && cfg.PageID != "",
	}
}

// GetName returns the publisher name
func (p *Publisher) GetName() string {
	return "Facebook"
}

// IsEnabled returns whether the publisher is enabled
func (p *Publisher) IsEnabled() bool {
	return p.enabled
}

// Publish publishes content to Facebook
func (p *Publisher) Publish(content *types.ProcessedContent) (*types.PublishResult, error) {
	result := &types.PublishResult{
		Platform:  types.PlatformFacebook,
		Timestamp: time.Now(),
	}

	if !p.enabled {
		result.Success = false
		result.Error = "Facebook publisher is not enabled or configured"
		return result, fmt.Errorf("publisher not enabled")
	}

	// Handle different content types
	switch content.Type {
	case types.ContentTypeText:
		return p.publishText(content, result)
	case types.ContentTypeImage:
		return p.publishWithImages(content, result)
	case types.ContentTypeVideo:
		return p.publishVideo(content, result)
	case types.ContentTypeMixed:
		// Facebook supports multiple images
		if len(content.MediaURLs) > 0 {
			return p.publishWithImages(content, result)
		}
		return p.publishText(content, result)
	default:
		result.Success = false
		result.Error = "unsupported content type"
		return result, fmt.Errorf("unsupported content type: %s", content.Type)
	}
}

// publishText publishes text-only post to Facebook
func (p *Publisher) publishText(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	// Facebook Graph API endpoint for page feed
	apiURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/feed", p.config.PageID)

	params := url.Values{}
	params.Set("message", content.Description)
	params.Set("access_token", p.config.AccessToken)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(params.Encode()))
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to send request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		result.Success = false
		result.Error = fmt.Sprintf("API error: %s (status: %d)", string(body), resp.StatusCode)
		return result, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var postResp struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &postResp); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to parse response: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = postResp.ID
	result.PostURL = fmt.Sprintf("https://facebook.com/%s", postResp.ID)

	return result, nil
}

// publishWithImages publishes post with images to Facebook
func (p *Publisher) publishWithImages(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	if len(content.MediaURLs) == 0 {
		return p.publishText(content, result)
	}

	// For single image, use simple photo post
	if len(content.MediaURLs) == 1 {
		return p.publishSinglePhoto(content, result)
	}

	// For multiple images, create album/carousel post
	return p.publishMultiplePhotos(content, result)
}

// publishSinglePhoto publishes single photo to Facebook
func (p *Publisher) publishSinglePhoto(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/photos", p.config.PageID)

	params := url.Values{}
	params.Set("url", content.MediaURLs[0])
	params.Set("caption", content.Description)
	params.Set("access_token", p.config.AccessToken)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(params.Encode()))
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to send request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		result.Success = false
		result.Error = fmt.Sprintf("API error: %s (status: %d)", string(body), resp.StatusCode)
		return result, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var photoResp struct {
		ID      string `json:"id"`
		PostID  string `json:"post_id"`
		Success bool   `json:"success"`
	}

	if err := json.Unmarshal(body, &photoResp); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to parse response: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = photoResp.PostID
	if result.PostID == "" {
		result.PostID = photoResp.ID
	}
	result.PostURL = fmt.Sprintf("https://facebook.com/%s", result.PostID)

	return result, nil
}

// publishMultiplePhotos publishes multiple photos to Facebook
func (p *Publisher) publishMultiplePhotos(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	// Step 1: Upload all photos unpublished
	var photoIDs []string
	for _, imageURL := range content.MediaURLs {
		photoID, err := p.uploadUnpublishedPhoto(imageURL)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("failed to upload photo: %v", err)
			return result, err
		}
		photoIDs = append(photoIDs, photoID)
	}

	// Step 2: Create post with all photos
	apiURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/feed", p.config.PageID)

	params := url.Values{}
	params.Set("message", content.Description)
	params.Set("access_token", p.config.AccessToken)

	// Attach all photos
	for _, photoID := range photoIDs {
		params.Add("attached_media[]", fmt.Sprintf(`{"media_fbid":"%s"}`, photoID))
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(params.Encode()))
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to send request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		result.Success = false
		result.Error = fmt.Sprintf("API error: %s (status: %d)", string(body), resp.StatusCode)
		return result, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var postResp struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &postResp); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to parse response: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = postResp.ID
	result.PostURL = fmt.Sprintf("https://facebook.com/%s", postResp.ID)

	return result, nil
}

// uploadUnpublishedPhoto uploads photo without publishing
func (p *Publisher) uploadUnpublishedPhoto(imageURL string) (string, error) {
	apiURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/photos", p.config.PageID)

	params := url.Values{}
	params.Set("url", imageURL)
	params.Set("published", "false")
	params.Set("access_token", p.config.AccessToken)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var photoResp struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &photoResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return photoResp.ID, nil
}

// publishVideo publishes video to Facebook
func (p *Publisher) publishVideo(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	if len(content.MediaURLs) == 0 {
		result.Success = false
		result.Error = "no video URL provided"
		return result, fmt.Errorf("no video URL")
	}

	videoURL := content.MediaURLs[0]

	// Download video
	videoData, err := p.downloadMedia(videoURL)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to download video: %v", err)
		return result, err
	}

	// Upload video to Facebook
	apiURL := fmt.Sprintf("https://graph-video.facebook.com/v18.0/%s/videos", p.config.PageID)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add video file
	part, err := writer.CreateFormFile("source", "video.mp4")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create form: %v", err)
		return result, err
	}
	part.Write(videoData)

	// Add description
	writer.WriteField("description", content.Description)
	writer.WriteField("access_token", p.config.AccessToken)

	writer.Close()

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := p.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to send request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		result.Success = false
		result.Error = fmt.Sprintf("API error: %s (status: %d)", string(respBody), resp.StatusCode)
		return result, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var videoResp struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(respBody, &videoResp); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to parse response: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = videoResp.ID
	result.PostURL = fmt.Sprintf("https://facebook.com/%s", videoResp.ID)

	return result, nil
}

// downloadMedia downloads media from URL
func (p *Publisher) downloadMedia(url string) ([]byte, error) {
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
		return nil, fmt.Errorf("failed to read media data: %w", err)
	}

	return data, nil
}
