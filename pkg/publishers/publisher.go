package publishers

import (
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
)

// Publisher interface for all social media publishers
type Publisher interface {
	// Publish publishes content to the platform
	Publish(content *types.ProcessedContent) (*types.PublishResult, error)

	// GetName returns the publisher name
	GetName() string

	// IsEnabled returns whether the publisher is enabled
	IsEnabled() bool
}
