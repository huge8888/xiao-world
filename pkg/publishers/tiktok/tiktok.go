package tiktok

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

// Publisher handles TikTok publishing
type Publisher struct {
	config     *configs.TikTokConfig
	httpClient *http.Client
	enabled    bool
}

// NewPublisher creates a new TikTok publisher
func NewPublisher(cfg *configs.TikTokConfig) *Publisher {
	return &Publisher{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for video uploads
		},
		enabled: cfg != nil && cfg.Enabled && cfg.AccessToken != "",
	}
}

// GetName returns the publisher name
func (p *Publisher) GetName() string {
	return "TikTok"
}

// IsEnabled returns whether the publisher is enabled
func (p *Publisher) IsEnabled() bool {
	return p.enabled
}

// Publish publishes content to TikTok
func (p *Publisher) Publish(content *types.ProcessedContent) (*types.PublishResult, error) {
	result := &types.PublishResult{
		Platform:  types.PlatformTikTok,
		Timestamp: time.Now(),
	}

	if !p.enabled {
		result.Success = false
		result.Error = "TikTok publisher is not enabled or configured"
		return result, fmt.Errorf("publisher not enabled")
	}

	// TikTok primarily supports video content
	switch content.Type {
	case types.ContentTypeVideo:
		return p.publishVideo(content, result)
	case types.ContentTypeImage:
		result.Success = false
		result.Error = "TikTok primarily supports video content. Image publishing not supported."
		return result, fmt.Errorf("image not supported")
	case types.ContentTypeText:
		result.Success = false
		result.Error = "TikTok requires video content. Text-only posts not supported."
		return result, fmt.Errorf("text-only not supported")
	default:
		result.Success = false
		result.Error = "unsupported content type"
		return result, fmt.Errorf("unsupported content type: %s", content.Type)
	}
}

// publishVideo publishes video to TikTok using Content Posting API
func (p *Publisher) publishVideo(content *types.ProcessedContent, result *types.PublishResult) (*types.PublishResult, error) {
	// TikTok Content Posting API workflow:
	// 1. Initialize video upload
	// 2. Upload video chunks
	// 3. Create post with video

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

	// Step 2: Initialize upload
	uploadURL, uploadID, err := p.initializeUpload()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to initialize upload: %v", err)
		return result, err
	}

	// Step 3: Upload video
	if err := p.uploadVideo(uploadURL, videoData); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to upload video: %v", err)
		return result, err
	}

	// Step 4: Create post
	postID, postURL, err := p.createPost(uploadID, content)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create post: %v", err)
		return result, err
	}

	result.Success = true
	result.PostID = postID
	result.PostURL = postURL

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

// initializeUpload initializes video upload to TikTok
func (p *Publisher) initializeUpload() (uploadURL string, uploadID string, err error) {
	// TikTok Content Posting API v2 endpoint
	apiURL := "https://open.tiktokapis.com/v2/post/publish/video/init/"

	reqBody := map[string]interface{}{
		"post_info": map[string]interface{}{
			"privacy_level": "PUBLIC_TO_EVERYONE",
		},
		"source_info": map[string]interface{}{
			"source": "FILE_UPLOAD",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.AccessToken))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var result struct {
		Data struct {
			UploadURL string `json:"upload_url"`
			PublishID string `json:"publish_id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data.UploadURL, result.Data.PublishID, nil
}

// uploadVideo uploads video data to TikTok
func (p *Publisher) uploadVideo(uploadURL string, videoData []byte) error {
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewBuffer(videoData))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Content-Type", "video/mp4")
	req.ContentLength = int64(len(videoData))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s (status: %d)", string(body), resp.StatusCode)
	}

	return nil
}

// createPost creates TikTok post with uploaded video
func (p *Publisher) createPost(publishID string, content *types.ProcessedContent) (string, string, error) {
	apiURL := "https://open.tiktokapis.com/v2/post/publish/status/fetch/"

	reqBody := map[string]interface{}{
		"publish_id": publishID,
	}

	// Add caption if available
	caption := content.Description
	if len(caption) > 2200 {
		caption = caption[:2200] // TikTok limit
	}

	if caption != "" {
		reqBody["post_info"] = map[string]interface{}{
			"title": caption,
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.AccessToken))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var result struct {
		Data struct {
			Status   string `json:"status"`
			ShareURL string `json:"share_url"`
			VideoID  string `json:"video_id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data.VideoID, result.Data.ShareURL, nil
}
