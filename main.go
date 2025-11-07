package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/processor"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/publishers"
	facebookPublisher "github.com/xpzouying/xiaohongshu-mcp/pkg/publishers/facebook"
	tiktokPublisher "github.com/xpzouying/xiaohongshu-mcp/pkg/publishers/tiktok"
	twitterPublisher "github.com/xpzouying/xiaohongshu-mcp/pkg/publishers/twitter"
	youtubePublisher "github.com/xpzouying/xiaohongshu-mcp/pkg/publishers/youtube"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/scheduler"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/translator"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
)

func main() {
	var (
		headless   bool
		binPath    string // 浏览器二进制文件路径
		port       string
		configPath string // 发布平台配置文件路径
	)
	flag.BoolVar(&headless, "headless", true, "是否无头模式")
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.StringVar(&port, "port", ":18060", "端口")
	flag.StringVar(&configPath, "config", "", "发布平台配置文件路径")
	flag.Parse()

	if len(binPath) == 0 {
		binPath = os.Getenv("ROD_BROWSER_BIN")
	}

	configs.InitHeadless(headless)
	configs.SetBinPath(binPath)

	// 初始化服务
	xiaohongshuService := NewXiaohongshuService()

	// 加载发布平台配置
	publishersConfig, err := configs.LoadPublishersConfig(configPath)
	if err != nil {
		logrus.Warnf("加载发布平台配置失败: %v，将使用环境变量", err)
		publishersConfig = configs.GetPublishersConfig()
	}

	// 初始化翻译器 - รองรับ AI หลายตัว (ChatGPT, Claude, Gemini) และ Google Translate
	var trans translator.Translator

	// ตรวจสอบว่าจะใช้ AI provider ไหน
	aiProvider := os.Getenv("AI_TRANSLATOR_PROVIDER") // openai, anthropic, google, google-translate
	aiAPIKey := os.Getenv("AI_TRANSLATOR_API_KEY")
	aiModel := os.Getenv("AI_TRANSLATOR_MODEL") // ไม่บังคับ จะใช้ค่าเริ่มต้น

	if aiProvider != "" && aiProvider != "google-translate" && aiAPIKey != "" {
		// ใช้ AI Translator
		trans = translator.NewAITranslator(aiProvider, aiAPIKey, aiModel)
		logrus.Infof("✅ ใช้ AI Translator: %s (model: %s)", aiProvider, aiModel)
	} else {
		// ใช้ Google Translate (เดิม)
		googleAPIKey := os.Getenv("GOOGLE_TRANSLATE_API_KEY")
		trans = translator.NewGoogleTranslator(googleAPIKey)
		if googleAPIKey == "" {
			logrus.Info("⚠️ ใช้ Google Translate ฟรี (มีข้อจำกัด rate limit)")
		} else {
			logrus.Info("✅ ใช้ Google Translate API")
		}
	}

	// 初始化内容处理器
	proc := processor.NewProcessor(trans)

	// 初始化各平台发布器
	publishersMap := make(map[types.Platform]publishers.Publisher)

	twitterPub := twitterPublisher.NewPublisher(publishersConfig.Twitter)
	if twitterPub.IsEnabled() {
		publishersMap[types.PlatformTwitter] = twitterPub
		logrus.Info("✅ Twitter 发布器已启用")
	} else {
		logrus.Info("⚠️ Twitter 发布器未启用")
	}

	tiktokPub := tiktokPublisher.NewPublisher(publishersConfig.TikTok)
	if tiktokPub.IsEnabled() {
		publishersMap[types.PlatformTikTok] = tiktokPub
		logrus.Info("✅ TikTok 发布器已启用")
	} else {
		logrus.Info("⚠️ TikTok 发布器未启用")
	}

	facebookPub := facebookPublisher.NewPublisher(publishersConfig.Facebook)
	if facebookPub.IsEnabled() {
		publishersMap[types.PlatformFacebook] = facebookPub
		logrus.Info("✅ Facebook 发布器已启用")
	} else {
		logrus.Info("⚠️ Facebook 发布器未启用")
	}

	youtubePub := youtubePublisher.NewPublisher(publishersConfig.YouTube)
	if youtubePub.IsEnabled() {
		publishersMap[types.PlatformYouTube] = youtubePub
		logrus.Info("✅ YouTube 发布器已启用")
	} else {
		logrus.Info("⚠️ YouTube 发布器未启用")
	}

	// 初始化调度器
	sched := scheduler.NewScheduler(proc, publishersMap)
	sched.Start()
	defer sched.Stop()

	logrus.Infof("🚀 多平台发布系统已初始化，启用了 %d 个平台", len(publishersMap))

	// 创建并启动应用服务器
	appServer := NewAppServer(xiaohongshuService)
	appServer.scheduler = sched
	appServer.publishers = publishersMap

	if err := appServer.Start(port); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
