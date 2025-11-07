package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"runtime/debug"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
)

// โครงสร้างพารามิเตอร์สำหรับ MCP tools

// PublishContentArgs พารามิเตอร์สำหรับเผยแพร่เนื้อหา
type PublishContentArgs struct {
	Title   string   `json:"title" jsonschema:"หัวข้อเนื้อหา (ข้อจำกัดของเสี้ยวหงชู: สูงสุด 20 คำภาษาจีนหรือคำภาษาอังกฤษ)"`
	Content string   `json:"content" jsonschema:"เนื้อหาหลัก ไม่รวม tags ที่ขึ้นต้นด้วย # ให้ใช้พารามิเตอร์ tags แทน"`
	Images  []string `json:"images" jsonschema:"รายการเส้นทางรูปภาพ (ต้องมีอย่างน้อย 1 รูป) รองรับ 2 วิธี: 1. ลิงก์ HTTP/HTTPS (ดาวน์โหลดอัตโนมัติ) 2. เส้นทางรูปภาพในเครื่อง (แนะนำ, เช่น: /Users/user/image.jpg)"`
	Tags    []string `json:"tags,omitempty" jsonschema:"รายการ tags หัวข้อ (ไม่บังคับ) เช่น [อาหาร, ท่องเที่ยว, ชีวิต]"`
}

// PublishVideoArgs พารามิเตอร์สำหรับเผยแพร่วิดีโอ (รองรับเฉพาะไฟล์วิดีโอในเครื่อง 1 ไฟล์)
type PublishVideoArgs struct {
	Title   string   `json:"title" jsonschema:"หัวข้อเนื้อหา (ข้อจำกัดของเสี้ยวหงชู: สูงสุด 20 คำภาษาจีนหรือคำภาษาอังกฤษ)"`
	Content string   `json:"content" jsonschema:"เนื้อหาหลัก ไม่รวม tags ที่ขึ้นต้นด้วย # ให้ใช้พารามิเตอร์ tags แทน"`
	Video   string   `json:"video" jsonschema:"เส้นทางวิดีโอในเครื่อง (รองรับเฉพาะไฟล์เดียว เช่น: /Users/user/video.mp4)"`
	Tags    []string `json:"tags,omitempty" jsonschema:"รายการ tags หัวข้อ (ไม่บังคับ) เช่น [อาหาร, ท่องเที่ยว, ชีวิต]"`
}

// SearchFeedsArgs พารามิเตอร์สำหรับค้นหาเนื้อหา
type SearchFeedsArgs struct {
	Keyword string       `json:"keyword" jsonschema:"คำค้นหา"`
	Filters FilterOption `json:"filters,omitempty" jsonschema:"ตัวเลือกกรอง"`
}

// FilterOption โครงสร้างตัวเลือกการกรอง
type FilterOption struct {
	SortBy      string `json:"sort_by,omitempty" jsonschema:"เรียงตาม: 综合|最新|最多点赞|最多评论|最多收藏, ค่าเริ่มต้น '综合'"`
	NoteType    string `json:"note_type,omitempty" jsonschema:"ประเภทโน้ต: 不限|视频|图文, ค่าเริ่มต้น '不限'"`
	PublishTime string `json:"publish_time,omitempty" jsonschema:"เวลาเผยแพร่: 不限|一天内|一周内|半年内, ค่าเริ่มต้น '不限'"`
	SearchScope string `json:"search_scope,omitempty" jsonschema:"ขอบเขตการค้นหา: 不限|已看过|未看过|已关注, ค่าเริ่มต้น '不限'"`
	Location    string `json:"location,omitempty" jsonschema:"ระยะตำแหน่ง: 不限|同城|附近, ค่าเริ่มต้น '不限'"`
}

// FeedDetailArgs พารามิเตอร์สำหรับดึงรายละเอียด Feed
type FeedDetailArgs struct {
	FeedID    string `json:"feed_id" jsonschema:"ID โน้ตเสี้ยวหงชู ดึงจากรายการ Feed"`
	XsecToken string `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
}

// UserProfileArgs พารามิเตอร์สำหรับดึงหน้าโปรไฟล์ผู้ใช้
type UserProfileArgs struct {
	UserID    string `json:"user_id" jsonschema:"ID ผู้ใช้เสี้ยวหงชู ดึงจากรายการ Feed"`
	XsecToken string `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
}

// PostCommentArgs พารามิเตอร์สำหรับแสดงความคิดเห็น
type PostCommentArgs struct {
	FeedID    string `json:"feed_id" jsonschema:"ID โน้ตเสี้ยวหงชู ดึงจากรายการ Feed"`
	XsecToken string `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
	Content   string `json:"content" jsonschema:"เนื้อหาความคิดเห็น"`
}

// LikeFeedArgs พารามิเตอร์สำหรับกดไลค์
type LikeFeedArgs struct {
	FeedID    string `json:"feed_id" jsonschema:"ID โน้ตเสี้ยวหงชู ดึงจากรายการ Feed"`
	XsecToken string `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
	Unlike    bool   `json:"unlike,omitempty" jsonschema:"ยกเลิกไลค์หรือไม่ true=ยกเลิกไลค์, false หรือไม่ระบุ=ไลค์"`
}

// FavoriteFeedArgs พารามิเตอร์สำหรับบันทึก
type FavoriteFeedArgs struct {
	FeedID     string `json:"feed_id" jsonschema:"ID โน้ตเสี้ยวหงชู ดึงจากรายการ Feed"`
	XsecToken  string `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
	Unfavorite bool   `json:"unfavorite,omitempty" jsonschema:"ยกเลิกบันทึกหรือไม่ true=ยกเลิกบันทึก, false หรือไม่ระบุ=บันทึก"`
}

// PublishToPlatformArgs พารามิเตอร์สำหรับเผยแพร่ไปแพลตฟอร์มเฉพาะ
type PublishToPlatformArgs struct {
	FeedID    string `json:"feed_id" jsonschema:"ID โน้ตเสี้ยวหงชู ดึงจากรายการ Feed หรือผลค้นหา"`
	XsecToken string `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
}

// PublishToAllPlatformsArgs พารามิเตอร์สำหรับเผยแพร่ไปทุกแพลตฟอร์ม
type PublishToAllPlatformsArgs struct {
	FeedID    string   `json:"feed_id" jsonschema:"ID โน้ตเสี้ยวหงชู ดึงจากรายการ Feed หรือผลค้นหา"`
	XsecToken string   `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
	Platforms []string `json:"platforms,omitempty" jsonschema:"รายการแพลตฟอร์ม (ไม่บังคับ) รองรับ: twitter, tiktok, facebook, youtube ถ้าไม่ระบุจะเผยแพร่ไปทุกแพลตฟอร์มที่เปิดใช้งาน"`
}

// SchedulePublishArgs พารามิเตอร์สำหรับกำหนดเวลาเผยแพร่
type SchedulePublishArgs struct {
	FeedID      string   `json:"feed_id" jsonschema:"ID โน้ตเสี้ยวหงชู ดึงจากรายการ Feed หรือผลค้นหา"`
	XsecToken   string   `json:"xsec_token" jsonschema:"Access token ดึงจากฟิลด์ xsecToken ในรายการ Feed"`
	Platforms   []string `json:"platforms" jsonschema:"รายการแพลตฟอร์ม รองรับ: twitter, tiktok, facebook, youtube"`
	ScheduledAt string   `json:"scheduled_at" jsonschema:"เวลาที่จะเผยแพร่ รูปแบบ: 2006-01-02 15:04:05"`
}

// CancelScheduledJobArgs พารามิเตอร์สำหรับยกเลิกงานที่กำหนดเวลา
type CancelScheduledJobArgs struct {
	JobID string `json:"job_id" jsonschema:"ID งาน ดึงจากผลลัพธ์ของ schedule_publish หรือ list_scheduled_jobs"`
}

// InitMCPServer 初始化 MCP Server
func InitMCPServer(appServer *AppServer) *mcp.Server {
	// 创建 MCP Server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "xiaohongshu-mcp",
			Version: "2.0.0",
		},
		nil,
	)

	// 注册所有工具
	registerTools(server, appServer)

	logrus.Info("MCP Server initialized with official SDK")

	return server
}

func withPanicRecovery[T any](
	toolName string,
	handler func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error),
) func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error) {

	return func(ctx context.Context, req *mcp.CallToolRequest, args T) (result *mcp.CallToolResult, resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"tool":  toolName,
					"panic": r,
				}).Error("Tool handler panicked")

				logrus.Errorf("Stack trace:\n%s", debug.Stack())

				result = &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("工具 %s 执行时发生内部错误: %v\n\n请查看服务端日志获取详细信息。", toolName, r),
						},
					},
					IsError: true,
				}
				resp = nil
				err = nil
			}
		}()

		return handler(ctx, req, args)
	}
}

// registerTools 注册所有 MCP 工具
func registerTools(server *mcp.Server, appServer *AppServer) {
	// 工具 1: 检查登录状态
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "check_login_status",
			Description: "检查小红书登录状态",
		},
		withPanicRecovery("check_login_status", func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			result := appServer.handleCheckLoginStatus(ctx)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 2: 获取登录二维码
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_login_qrcode",
			Description: "获取登录二维码（返回 Base64 图片和超时时间）",
		},
		withPanicRecovery("get_login_qrcode", func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			result := appServer.handleGetLoginQrcode(ctx)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 3: 删除 cookies（登录重置）
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "delete_cookies",
			Description: "删除 cookies 文件，重置登录状态。删除后需要重新登录。",
		},
		withPanicRecovery("delete_cookies", func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			result := appServer.handleDeleteCookies(ctx)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 4: 发布内容
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_content",
			Description: "发布小红书图文内容",
		},
		withPanicRecovery("publish_content", func(ctx context.Context, req *mcp.CallToolRequest, args PublishContentArgs) (*mcp.CallToolResult, any, error) {
			// 转换参数格式到现有的 handler
			argsMap := map[string]interface{}{
				"title":   args.Title,
				"content": args.Content,
				"images":  convertStringsToInterfaces(args.Images),
				"tags":    convertStringsToInterfaces(args.Tags),
			}
			result := appServer.handlePublishContent(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 5: 获取Feed列表
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "list_feeds",
			Description: "获取首页 Feeds 列表",
		},
		withPanicRecovery("list_feeds", func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			result := appServer.handleListFeeds(ctx)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 6: 搜索内容
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "search_feeds",
			Description: "搜索小红书内容（需要已登录）",
		},
		withPanicRecovery("search_feeds", func(ctx context.Context, req *mcp.CallToolRequest, args SearchFeedsArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handleSearchFeeds(ctx, args)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 7: 获取Feed详情
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_feed_detail",
			Description: "获取小红书笔记详情，返回笔记内容、图片、作者信息、互动数据（点赞/收藏/分享数）及评论列表",
		},
		withPanicRecovery("get_feed_detail", func(ctx context.Context, req *mcp.CallToolRequest, args FeedDetailArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
			}
			result := appServer.handleGetFeedDetail(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 8: 获取用户主页
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "user_profile",
			Description: "获取指定的小红书用户主页，返回用户基本信息，关注、粉丝、获赞量及其笔记内容",
		},
		withPanicRecovery("user_profile", func(ctx context.Context, req *mcp.CallToolRequest, args UserProfileArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"user_id":    args.UserID,
				"xsec_token": args.XsecToken,
			}
			result := appServer.handleUserProfile(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 9: 发表评论
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "post_comment_to_feed",
			Description: "发表评论到小红书笔记",
		},
		withPanicRecovery("post_comment_to_feed", func(ctx context.Context, req *mcp.CallToolRequest, args PostCommentArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
				"content":    args.Content,
			}
			result := appServer.handlePostComment(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 10: 发布视频（仅本地文件）
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_with_video",
			Description: "发布小红书视频内容（仅支持本地单个视频文件）",
		},
		withPanicRecovery("publish_with_video", func(ctx context.Context, req *mcp.CallToolRequest, args PublishVideoArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"title":   args.Title,
				"content": args.Content,
				"video":   args.Video,
				"tags":    convertStringsToInterfaces(args.Tags),
			}
			result := appServer.handlePublishVideo(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 11: 点赞笔记
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "like_feed",
			Description: "为指定笔记点赞或取消点赞（如已点赞将跳过点赞，如未点赞将跳过取消点赞）",
		},
		withPanicRecovery("like_feed", func(ctx context.Context, req *mcp.CallToolRequest, args LikeFeedArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
				"unlike":     args.Unlike,
			}
			result := appServer.handleLikeFeed(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 12: 收藏笔记
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "favorite_feed",
			Description: "收藏指定笔记或取消收藏（如已收藏将跳过收藏，如未收藏将跳过取消收藏）",
		},
		withPanicRecovery("favorite_feed", func(ctx context.Context, req *mcp.CallToolRequest, args FavoriteFeedArgs) (*mcp.CallToolResult, any, error) {
			argsMap := map[string]interface{}{
				"feed_id":    args.FeedID,
				"xsec_token": args.XsecToken,
				"unfavorite": args.Unfavorite,
			}
			result := appServer.handleFavoriteFeed(ctx, argsMap)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 13: 发布到 Twitter
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_to_twitter",
			Description: "将小红书笔记内容发布到 Twitter/X（自动翻译为英文）",
		},
		withPanicRecovery("publish_to_twitter", func(ctx context.Context, req *mcp.CallToolRequest, args PublishToPlatformArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handlePublishToPlatform(ctx, args.FeedID, args.XsecToken, "twitter")
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 14: 发布到 TikTok
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_to_tiktok",
			Description: "将小红书笔记内容发布到 TikTok（仅支持视频内容，自动翻译为英文）",
		},
		withPanicRecovery("publish_to_tiktok", func(ctx context.Context, req *mcp.CallToolRequest, args PublishToPlatformArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handlePublishToPlatform(ctx, args.FeedID, args.XsecToken, "tiktok")
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 15: 发布到 Facebook
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_to_facebook",
			Description: "将小红书笔记内容发布到 Facebook（自动翻译为英文）",
		},
		withPanicRecovery("publish_to_facebook", func(ctx context.Context, req *mcp.CallToolRequest, args PublishToPlatformArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handlePublishToPlatform(ctx, args.FeedID, args.XsecToken, "facebook")
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 16: 发布到 YouTube
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_to_youtube",
			Description: "将小红书笔记内容发布到 YouTube（仅支持视频内容，自动翻译为英文）",
		},
		withPanicRecovery("publish_to_youtube", func(ctx context.Context, req *mcp.CallToolRequest, args PublishToPlatformArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handlePublishToPlatform(ctx, args.FeedID, args.XsecToken, "youtube")
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 17: 发布到所有平台
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "publish_to_all_platforms",
			Description: "将小红书笔记内容同时发布到多个平台（Twitter, TikTok, Facebook, YouTube），自动翻译为英文",
		},
		withPanicRecovery("publish_to_all_platforms", func(ctx context.Context, req *mcp.CallToolRequest, args PublishToAllPlatformsArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handlePublishToAllPlatforms(ctx, args)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 18: 定时发布
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "schedule_publish",
			Description: "创建定时发布任务，在指定时间将小红书笔记发布到指定平台",
		},
		withPanicRecovery("schedule_publish", func(ctx context.Context, req *mcp.CallToolRequest, args SchedulePublishArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handleSchedulePublish(ctx, args)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 19: 查看定时任务列表
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "list_scheduled_jobs",
			Description: "查看所有定时发布任务及其状态",
		},
		withPanicRecovery("list_scheduled_jobs", func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
			result := appServer.handleListScheduledJobs(ctx)
			return convertToMCPResult(result), nil, nil
		}),
	)

	// 工具 20: 取消定时任务
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "cancel_scheduled_job",
			Description: "取消指定的定时发布任务",
		},
		withPanicRecovery("cancel_scheduled_job", func(ctx context.Context, req *mcp.CallToolRequest, args CancelScheduledJobArgs) (*mcp.CallToolResult, any, error) {
			result := appServer.handleCancelScheduledJob(ctx, args.JobID)
			return convertToMCPResult(result), nil, nil
		}),
	)

	logrus.Infof("Registered %d MCP tools", 20)
}

// convertToMCPResult 将自定义的 MCPToolResult 转换为官方 SDK 的格式
func convertToMCPResult(result *MCPToolResult) *mcp.CallToolResult {
	var contents []mcp.Content
	for _, c := range result.Content {
		switch c.Type {
		case "text":
			contents = append(contents, &mcp.TextContent{Text: c.Text})
		case "image":
			// 解码 base64 字符串为 []byte
			imageData, err := base64.StdEncoding.DecodeString(c.Data)
			if err != nil {
				logrus.WithError(err).Error("Failed to decode base64 image data")
				// 如果解码失败，添加错误文本
				contents = append(contents, &mcp.TextContent{
					Text: "图片数据解码失败: " + err.Error(),
				})
			} else {
				contents = append(contents, &mcp.ImageContent{
					Data:     imageData,
					MIMEType: c.MimeType,
				})
			}
		}
	}

	return &mcp.CallToolResult{
		Content: contents,
		IsError: result.IsError,
	}
}

// convertStringsToInterfaces 辅助函数：将 []string 转换为 []interface{}
func convertStringsToInterfaces(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
