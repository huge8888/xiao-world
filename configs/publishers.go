package configs

import (
	"encoding/json"
	"os"
)

// PublishersConfig holds configuration for all publishing platforms
type PublishersConfig struct {
	Twitter  *TwitterConfig  `json:"twitter"`
	TikTok   *TikTokConfig   `json:"tiktok"`
	Facebook *FacebookConfig `json:"facebook"`
	YouTube  *YouTubeConfig  `json:"youtube"`
}

// TwitterConfig holds Twitter/X API configuration
type TwitterConfig struct {
	Enabled     bool   `json:"enabled"`
	BearerToken string `json:"bearer_token"` // OAuth 2.0 Bearer Token
	APIKey      string `json:"api_key"`      // Optional: for OAuth 1.0a
	APISecret   string `json:"api_secret"`   // Optional: for OAuth 1.0a
}

// TikTokConfig holds TikTok API configuration
type TikTokConfig struct {
	Enabled      bool   `json:"enabled"`
	AccessToken  string `json:"access_token"`  // TikTok Content Posting API access token
	ClientKey    string `json:"client_key"`    // TikTok app client key
	ClientSecret string `json:"client_secret"` // TikTok app client secret
}

// FacebookConfig holds Facebook API configuration
type FacebookConfig struct {
	Enabled     bool   `json:"enabled"`
	AccessToken string `json:"access_token"` // Facebook Page Access Token
	PageID      string `json:"page_id"`      // Facebook Page ID to post to
}

// YouTubeConfig holds YouTube API configuration
type YouTubeConfig struct {
	Enabled      bool   `json:"enabled"`
	AccessToken  string `json:"access_token"`  // OAuth 2.0 Access Token
	RefreshToken string `json:"refresh_token"` // OAuth 2.0 Refresh Token
	ClientID     string `json:"client_id"`     // OAuth 2.0 Client ID
	ClientSecret string `json:"client_secret"` // OAuth 2.0 Client Secret
}

var globalPublishersConfig *PublishersConfig

// LoadPublishersConfig loads publishers configuration from file or environment variables
func LoadPublishersConfig(configPath string) (*PublishersConfig, error) {
	config := &PublishersConfig{
		Twitter:  loadTwitterConfig(configPath),
		TikTok:   loadTikTokConfig(configPath),
		Facebook: loadFacebookConfig(configPath),
		YouTube:  loadYouTubeConfig(configPath),
	}

	globalPublishersConfig = config
	return config, nil
}

// GetPublishersConfig returns the global publishers configuration
func GetPublishersConfig() *PublishersConfig {
	if globalPublishersConfig == nil {
		// Load default config from environment variables
		globalPublishersConfig = &PublishersConfig{
			Twitter:  loadTwitterConfig(""),
			TikTok:   loadTikTokConfig(""),
			Facebook: loadFacebookConfig(""),
			YouTube:  loadYouTubeConfig(""),
		}
	}
	return globalPublishersConfig
}

// loadTwitterConfig loads Twitter configuration
func loadTwitterConfig(configPath string) *TwitterConfig {
	config := &TwitterConfig{}

	// Try to load from file
	if configPath != "" {
		data, err := os.ReadFile(configPath + "/twitter.json")
		if err == nil {
			if err := json.Unmarshal(data, config); err == nil {
				return config
			}
		}
	}

	// Load from environment variables
	config.Enabled = os.Getenv("TWITTER_ENABLED") == "true"
	config.BearerToken = os.Getenv("TWITTER_BEARER_TOKEN")
	config.APIKey = os.Getenv("TWITTER_API_KEY")
	config.APISecret = os.Getenv("TWITTER_API_SECRET")

	return config
}

// loadTikTokConfig loads TikTok configuration
func loadTikTokConfig(configPath string) *TikTokConfig {
	config := &TikTokConfig{}

	// Try to load from file
	if configPath != "" {
		data, err := os.ReadFile(configPath + "/tiktok.json")
		if err == nil {
			if err := json.Unmarshal(data, config); err == nil {
				return config
			}
		}
	}

	// Load from environment variables
	config.Enabled = os.Getenv("TIKTOK_ENABLED") == "true"
	config.AccessToken = os.Getenv("TIKTOK_ACCESS_TOKEN")
	config.ClientKey = os.Getenv("TIKTOK_CLIENT_KEY")
	config.ClientSecret = os.Getenv("TIKTOK_CLIENT_SECRET")

	return config
}

// loadFacebookConfig loads Facebook configuration
func loadFacebookConfig(configPath string) *FacebookConfig {
	config := &FacebookConfig{}

	// Try to load from file
	if configPath != "" {
		data, err := os.ReadFile(configPath + "/facebook.json")
		if err == nil {
			if err := json.Unmarshal(data, config); err == nil {
				return config
			}
		}
	}

	// Load from environment variables
	config.Enabled = os.Getenv("FACEBOOK_ENABLED") == "true"
	config.AccessToken = os.Getenv("FACEBOOK_ACCESS_TOKEN")
	config.PageID = os.Getenv("FACEBOOK_PAGE_ID")

	return config
}

// loadYouTubeConfig loads YouTube configuration
func loadYouTubeConfig(configPath string) *YouTubeConfig {
	config := &YouTubeConfig{}

	// Try to load from file
	if configPath != "" {
		data, err := os.ReadFile(configPath + "/youtube.json")
		if err == nil {
			if err := json.Unmarshal(data, config); err == nil {
				return config
			}
		}
	}

	// Load from environment variables
	config.Enabled = os.Getenv("YOUTUBE_ENABLED") == "true"
	config.AccessToken = os.Getenv("YOUTUBE_ACCESS_TOKEN")
	config.RefreshToken = os.Getenv("YOUTUBE_REFRESH_TOKEN")
	config.ClientID = os.Getenv("YOUTUBE_CLIENT_ID")
	config.ClientSecret = os.Getenv("YOUTUBE_CLIENT_SECRET")

	return config
}
