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
		binPath    string // เส้นทางไฟล์ binary ของเบราว์เซอร์
		port       string
		configPath string // เส้นทางไฟล์ config ของแพลตฟอร์ม
	)
	flag.BoolVar(&headless, "headless", true, "ใช้โหมด headless หรือไม่")
	flag.StringVar(&binPath, "bin", "", "เส้นทางไฟล์ binary ของเบราว์เซอร์")
	flag.StringVar(&port, "port", ":18060", "พอร์ต")
	flag.StringVar(&configPath, "config", "", "เส้นทางไฟล์ config ของแพลตฟอร์ม")
	flag.Parse()

	if len(binPath) == 0 {
		binPath = os.Getenv("ROD_BROWSER_BIN")
	}

	configs.InitHeadless(headless)
	configs.SetBinPath(binPath)

	// เริ่มต้นบริการ
	xiaohongshuService := NewXiaohongshuService()

	// โหลดการตั้งค่าแพลตฟอร์ม
	publishersConfig, err := configs.LoadPublishersConfig(configPath)
	if err != nil {
		logrus.Warnf("โหลดการตั้งค่าแพลตฟอร์มล้มเหลว: %v, จะใช้ environment variables แทน", err)
		publishersConfig = configs.GetPublishersConfig()
	}

	// เริ่มต้น translator - รองรับ AI หลายตัว (ChatGPT, Claude, Gemini) และ Google Translate
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

	// เริ่มต้นตัวประมวลผลเนื้อหา
	proc := processor.NewProcessor(trans)

	// เริ่มต้น publisher แต่ละแพลตฟอร์ม
	publishersMap := make(map[types.Platform]publishers.Publisher)

	twitterPub := twitterPublisher.NewPublisher(publishersConfig.Twitter)
	if twitterPub.IsEnabled() {
		publishersMap[types.PlatformTwitter] = twitterPub
		logrus.Info("✅ Twitter publisher เปิดใช้งานแล้ว")
	} else {
		logrus.Info("⚠️ Twitter publisher ไม่ได้เปิดใช้งาน")
	}

	tiktokPub := tiktokPublisher.NewPublisher(publishersConfig.TikTok)
	if tiktokPub.IsEnabled() {
		publishersMap[types.PlatformTikTok] = tiktokPub
		logrus.Info("✅ TikTok publisher เปิดใช้งานแล้ว")
	} else {
		logrus.Info("⚠️ TikTok publisher ไม่ได้เปิดใช้งาน")
	}

	facebookPub := facebookPublisher.NewPublisher(publishersConfig.Facebook)
	if facebookPub.IsEnabled() {
		publishersMap[types.PlatformFacebook] = facebookPub
		logrus.Info("✅ Facebook publisher เปิดใช้งานแล้ว")
	} else {
		logrus.Info("⚠️ Facebook publisher ไม่ได้เปิดใช้งาน")
	}

	youtubePub := youtubePublisher.NewPublisher(publishersConfig.YouTube)
	if youtubePub.IsEnabled() {
		publishersMap[types.PlatformYouTube] = youtubePub
		logrus.Info("✅ YouTube publisher เปิดใช้งานแล้ว")
	} else {
		logrus.Info("⚠️ YouTube publisher ไม่ได้เปิดใช้งาน")
	}

	// เริ่มต้น scheduler
	sched := scheduler.NewScheduler(proc, publishersMap)
	sched.Start()
	defer sched.Stop()

	logrus.Infof("🚀 ระบบเผยแพร่หลายแพลตฟอร์มเริ่มต้นแล้ว เปิดใช้งาน %d แพลตฟอร์ม", len(publishersMap))

	// สร้างและเริ่มต้น app server
	appServer := NewAppServer(xiaohongshuService)
	appServer.scheduler = sched
	appServer.publishers = publishersMap

	if err := appServer.Start(port); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
