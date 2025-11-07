package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

// handlePublishToPlatform handles publishing to a specific platform
func (s *AppServer) handlePublishToPlatform(ctx context.Context, feedID, xsecToken, platformName string) *MCPToolResult {
	if s.scheduler == nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 多平台发布服务未初始化，请检查配置"},
			},
			IsError: true,
		}
	}

	// Get feed detail
	feedDetailResp, err := s.xiaohongshuService.GetFeedDetail(ctx, feedID, xsecToken)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 获取笔记详情失败: %v", err)},
			},
			IsError: true,
		}
	}

	// Extract actual FeedDetail from response
	feedDetail, ok := feedDetailResp.Data.(*xiaohongshu.FeedDetail)
	if !ok {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 获取笔记详情失败: 数据格式错误"},
			},
			IsError: true,
		}
	}

	// Convert platform name to Platform type
	var platform types.Platform
	switch platformName {
	case "twitter":
		platform = types.PlatformTwitter
	case "tiktok":
		platform = types.PlatformTikTok
	case "facebook":
		platform = types.PlatformFacebook
	case "youtube":
		platform = types.PlatformYouTube
	default:
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 不支持的平台: %s", platformName)},
			},
			IsError: true,
		}
	}

	// Publish immediately
	results, err := s.scheduler.PublishNow(feedDetail, []types.Platform{platform})
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 发布失败: %v", err)},
			},
			IsError: true,
		}
	}

	// Format results
	if len(results) == 0 {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 没有返回发布结果"},
			},
			IsError: true,
		}
	}

	result := results[0]
	if result.Success {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("✅ 成功发布到 %s\n\n📝 帖子ID: %s\n🔗 链接: %s\n⏰ 时间: %s",
					platformName, result.PostID, result.PostURL, result.Timestamp.Format("2006-01-02 15:04:05"))},
			},
			IsError: false,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: fmt.Sprintf("❌ 发布到 %s 失败: %s", platformName, result.Error)},
		},
		IsError: true,
	}
}

// handlePublishToAllPlatforms handles publishing to multiple platforms
func (s *AppServer) handlePublishToAllPlatforms(ctx context.Context, args PublishToAllPlatformsArgs) *MCPToolResult {
	if s.scheduler == nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 多平台发布服务未初始化，请检查配置"},
			},
			IsError: true,
		}
	}

	// Get feed detail
	feedDetailResp, err := s.xiaohongshuService.GetFeedDetail(ctx, args.FeedID, args.XsecToken)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 获取笔记详情失败: %v", err)},
			},
			IsError: true,
		}
	}

	// Extract actual FeedDetail from response
	feedDetail, ok := feedDetailResp.Data.(*xiaohongshu.FeedDetail)
	if !ok {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 获取笔记详情失败: 数据格式错误"},
			},
			IsError: true,
		}
	}

	// Determine platforms
	var platforms []types.Platform
	if len(args.Platforms) > 0 {
		// Use specified platforms
		for _, platformName := range args.Platforms {
			switch platformName {
			case "twitter":
				platforms = append(platforms, types.PlatformTwitter)
			case "tiktok":
				platforms = append(platforms, types.PlatformTikTok)
			case "facebook":
				platforms = append(platforms, types.PlatformFacebook)
			case "youtube":
				platforms = append(platforms, types.PlatformYouTube)
			default:
				logrus.Warnf("Unknown platform: %s", platformName)
			}
		}
	} else {
		// Use all enabled platforms
		for platform, publisher := range s.publishers {
			if publisher.IsEnabled() {
				platforms = append(platforms, platform)
			}
		}
	}

	if len(platforms) == 0 {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 没有可用的平台，请检查配置"},
			},
			IsError: true,
		}
	}

	// Publish to all platforms
	results, err := s.scheduler.PublishNow(feedDetail, platforms)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 发布失败: %v", err)},
			},
			IsError: true,
		}
	}

	// Format results
	var successCount int
	var failCount int
	resultText := "📊 发布结果汇总:\n\n"

	for _, result := range results {
		if result.Success {
			successCount++
			resultText += fmt.Sprintf("✅ %s: 成功\n   🔗 %s\n", result.Platform, result.PostURL)
		} else {
			failCount++
			resultText += fmt.Sprintf("❌ %s: 失败 - %s\n", result.Platform, result.Error)
		}
	}

	resultText += fmt.Sprintf("\n📈 总计: %d 成功, %d 失败", successCount, failCount)

	return &MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: resultText},
		},
		IsError: failCount > 0 && successCount == 0,
	}
}

// handleSchedulePublish handles scheduling a publish job
func (s *AppServer) handleSchedulePublish(ctx context.Context, args SchedulePublishArgs) *MCPToolResult {
	if s.scheduler == nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 多平台发布服务未初始化，请检查配置"},
			},
			IsError: true,
		}
	}

	// Parse scheduled time
	scheduledAt, err := time.Parse("2006-01-02 15:04:05", args.ScheduledAt)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 时间格式错误: %v\n正确格式: 2006-01-02 15:04:05", err)},
			},
			IsError: true,
		}
	}

	// Convert platform names to Platform types
	var platforms []types.Platform
	for _, platformName := range args.Platforms {
		switch platformName {
		case "twitter":
			platforms = append(platforms, types.PlatformTwitter)
		case "tiktok":
			platforms = append(platforms, types.PlatformTikTok)
		case "facebook":
			platforms = append(platforms, types.PlatformFacebook)
		case "youtube":
			platforms = append(platforms, types.PlatformYouTube)
		default:
			return &MCPToolResult{
				Content: []MCPContent{
					{Type: "text", Text: fmt.Sprintf("❌ 不支持的平台: %s", platformName)},
				},
				IsError: true,
			}
		}
	}

	// Schedule the job
	jobID, err := s.scheduler.ScheduleJob(args.FeedID, platforms, scheduledAt)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 创建定时任务失败: %v", err)},
			},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: fmt.Sprintf("✅ 定时任务创建成功\n\n🆔 任务ID: %s\n📅 发布时间: %s\n📱 平台: %v",
				jobID, scheduledAt.Format("2006-01-02 15:04:05"), args.Platforms)},
		},
		IsError: false,
	}
}

// handleListScheduledJobs handles listing all scheduled jobs
func (s *AppServer) handleListScheduledJobs(ctx context.Context) *MCPToolResult {
	if s.scheduler == nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 多平台发布服务未初始化，请检查配置"},
			},
			IsError: true,
		}
	}

	jobs := s.scheduler.ListJobs()

	if len(jobs) == 0 {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "📋 当前没有定时任务"},
			},
			IsError: false,
		}
	}

	// Format jobs list
	jobsJSON, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 格式化任务列表失败: %v", err)},
			},
			IsError: true,
		}
	}

	resultText := fmt.Sprintf("📋 定时任务列表 (共 %d 个):\n\n```json\n%s\n```", len(jobs), string(jobsJSON))

	return &MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: resultText},
		},
		IsError: false,
	}
}

// handleCancelScheduledJob handles canceling a scheduled job
func (s *AppServer) handleCancelScheduledJob(ctx context.Context, jobID string) *MCPToolResult {
	if s.scheduler == nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: "❌ 多平台发布服务未初始化，请检查配置"},
			},
			IsError: true,
		}
	}

	err := s.scheduler.CancelJob(jobID)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{Type: "text", Text: fmt.Sprintf("❌ 取消任务失败: %v", err)},
			},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: fmt.Sprintf("✅ 任务已取消\n\n🆔 任务ID: %s", jobID)},
		},
		IsError: false,
	}
}
