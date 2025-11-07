package types

import "time"

// Platform represents social media platforms
type Platform string

const (
	PlatformTwitter  Platform = "twitter"
	PlatformTikTok   Platform = "tiktok"
	PlatformFacebook Platform = "facebook"
	PlatformYouTube  Platform = "youtube"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
	ContentTypeVideo ContentType = "video"
	ContentTypeMixed ContentType = "mixed"
)

// ProcessedContent represents content after translation and adaptation
type ProcessedContent struct {
	Platform    Platform    `json:"platform"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Type        ContentType `json:"type"`
	MediaURLs   []string    `json:"media_urls"`
	Tags        []string    `json:"tags"`

	// Original info
	OriginalTitle       string `json:"original_title"`
	OriginalDescription string `json:"original_description"`
	SourceID            string `json:"source_id"`
	SourceURL           string `json:"source_url"`
}

// PublishRequest represents a request to publish content
type PublishRequest struct {
	FeedID     string     `json:"feed_id"`
	XsecToken  string     `json:"xsec_token"`
	Platforms  []Platform `json:"platforms"`
	ScheduleAt *time.Time `json:"schedule_at,omitempty"`
}

// PublishResult represents the result of a publish operation
type PublishResult struct {
	Platform  Platform  `json:"platform"`
	Success   bool      `json:"success"`
	PostID    string    `json:"post_id,omitempty"`
	PostURL   string    `json:"post_url,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// ScheduledJob represents a scheduled publishing job
type ScheduledJob struct {
	ID          string          `json:"id"`
	FeedID      string          `json:"feed_id"`
	XsecToken   string          `json:"xsec_token"`
	Platforms   []Platform      `json:"platforms"`
	ScheduledAt time.Time       `json:"scheduled_at"`
	Status      JobStatus       `json:"status"`
	Results     []PublishResult `json:"results,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Error       string          `json:"error,omitempty"`
}

// JobStatus represents job status
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)
