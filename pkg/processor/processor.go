package processor

import (
	"fmt"
	"strings"

	"github.com/xpzouying/xiaohongshu-mcp/pkg/translator"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

// Processor handles content processing and adaptation
type Processor struct {
	translator translator.Translator
}

// NewProcessor creates a new content processor
func NewProcessor(trans translator.Translator) *Processor {
	return &Processor{
		translator: trans,
	}
}

// Process processes Xiaohongshu content for a specific platform
func (p *Processor) Process(feed *xiaohongshu.FeedDetail, platform types.Platform) (*types.ProcessedContent, error) {
	// Translate content to English
	translatedTitle, err := p.translator.Translate(feed.Title, "zh", "en")
	if err != nil {
		return nil, fmt.Errorf("failed to translate title: %w", err)
	}

	translatedDesc, err := p.translator.Translate(feed.Desc, "zh", "en")
	if err != nil {
		return nil, fmt.Errorf("failed to translate description: %w", err)
	}

	// Determine content type
	contentType := types.ContentTypeText
	var mediaURLs []string

	if feed.Type == "video" {
		contentType = types.ContentTypeVideo
		// Note: Video URL needs to be extracted from browser or downloaded
		// For now, we'll mark it as video type
	} else if len(feed.ImageList) > 0 {
		contentType = types.ContentTypeImage
		for _, img := range feed.ImageList {
			if img.URLDefault != "" {
				mediaURLs = append(mediaURLs, img.URLDefault)
			}
		}
	}

	processed := &types.ProcessedContent{
		Platform:            platform,
		Title:               translatedTitle,
		Description:         translatedDesc,
		Type:                contentType,
		MediaURLs:           mediaURLs,
		Tags:                []string{},
		OriginalTitle:       feed.Title,
		OriginalDescription: feed.Desc,
		SourceID:            feed.NoteID,
		SourceURL:           fmt.Sprintf("https://www.xiaohongshu.com/explore/%s", feed.NoteID),
	}

	// Adapt content for specific platform
	switch platform {
	case types.PlatformTwitter:
		return p.adaptForTwitter(processed)
	case types.PlatformTikTok:
		return p.adaptForTikTok(processed)
	case types.PlatformFacebook:
		return p.adaptForFacebook(processed)
	case types.PlatformYouTube:
		return p.adaptForYouTube(processed)
	default:
		return processed, nil
	}
}

// adaptForTwitter adapts content for Twitter/X
func (p *Processor) adaptForTwitter(content *types.ProcessedContent) (*types.ProcessedContent, error) {
	// Twitter limit: 280 characters for text (4000 for Twitter Blue/Premium)
	maxLength := 280

	// Create tweet text with title and description
	tweetText := content.Title
	if content.Description != "" {
		// Add description if it fits
		combined := fmt.Sprintf("%s\n\n%s", content.Title, content.Description)
		if len(combined) <= maxLength {
			tweetText = combined
		} else {
			// Truncate description to fit
			remaining := maxLength - len(content.Title) - 5 // 5 for "\n\n" and "..."
			if remaining > 0 {
				truncated := p.truncateText(content.Description, remaining)
				tweetText = fmt.Sprintf("%s\n\n%s...", content.Title, truncated)
			}
		}
	}

	// Add source link if space allows
	sourceTag := fmt.Sprintf("\n\nSource: %s", content.SourceURL)
	if len(tweetText)+len(sourceTag) <= maxLength {
		tweetText += sourceTag
	}

	content.Description = tweetText

	// Twitter supports up to 4 images or 1 video
	if content.Type == types.ContentTypeImage && len(content.MediaURLs) > 4 {
		content.MediaURLs = content.MediaURLs[:4]
	}

	return content, nil
}

// adaptForTikTok adapts content for TikTok
func (p *Processor) adaptForTikTok(content *types.ProcessedContent) (*types.ProcessedContent, error) {
	// TikTok: video only, description up to 2200 characters
	maxLength := 2200

	// TikTok is video-focused
	if content.Type != types.ContentTypeVideo {
		return nil, fmt.Errorf("TikTok requires video content, got: %s", content.Type)
	}

	// Create description
	description := content.Title
	if content.Description != "" {
		combined := fmt.Sprintf("%s\n\n%s", content.Title, content.Description)
		if len(combined) <= maxLength {
			description = combined
		} else {
			truncated := p.truncateText(content.Description, maxLength-len(content.Title)-5)
			description = fmt.Sprintf("%s\n\n%s...", content.Title, truncated)
		}
	}

	// Add source attribution
	sourceTag := fmt.Sprintf("\n\n📱 From Xiaohongshu: %s", content.SourceURL)
	if len(description)+len(sourceTag) <= maxLength {
		description += sourceTag
	}

	content.Description = description

	return content, nil
}

// adaptForFacebook adapts content for Facebook
func (p *Processor) adaptForFacebook(content *types.ProcessedContent) (*types.ProcessedContent, error) {
	// Facebook: virtually no text limit, but optimal is 40-80 characters for engagement
	// Supports multiple images and videos

	// Create post text
	postText := fmt.Sprintf("%s\n\n%s", content.Title, content.Description)

	// Add source attribution
	postText += fmt.Sprintf("\n\n🔗 Original post from Xiaohongshu:\n%s", content.SourceURL)

	content.Description = postText

	return content, nil
}

// adaptForYouTube adapts content for YouTube
func (p *Processor) adaptForYouTube(content *types.ProcessedContent) (*types.ProcessedContent, error) {
	// YouTube: video only, title up to 100 characters, description up to 5000 characters
	maxTitleLength := 100
	maxDescLength := 5000

	// YouTube is video-only
	if content.Type != types.ContentTypeVideo {
		return nil, fmt.Errorf("YouTube requires video content, got: %s", content.Type)
	}

	// Truncate title if needed
	if len(content.Title) > maxTitleLength {
		content.Title = p.truncateText(content.Title, maxTitleLength-3) + "..."
	}

	// Create description
	description := content.Description
	if len(description) > maxDescLength {
		description = p.truncateText(description, maxDescLength-3) + "..."
	}

	// Add source attribution and metadata
	sourceInfo := fmt.Sprintf("\n\n━━━━━━━━━━━━━━━━━━━━━━\n📱 Original Content from Xiaohongshu\n🔗 Source: %s\n\n#Xiaohongshu #ContentSharing", content.SourceURL)
	if len(description)+len(sourceInfo) <= maxDescLength {
		description += sourceInfo
	}

	content.Description = description

	return content, nil
}

// truncateText truncates text to the specified length at word boundary
func (p *Processor) truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	// Try to truncate at word boundary
	truncated := text[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = text[:lastSpace]
	}

	return truncated
}
