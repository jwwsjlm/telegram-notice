package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	Tbot "telegram-notice/Newtele"
	"telegram-notice/global"
	"telegram-notice/hash"
	"telegram-notice/router"
)

func main() {
	// 设置 Gin 的模式为测试模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化日志
	InitLogger()
	defer global.LogZap.Sync() // 程序结束前同步日志缓冲区

	// 加载配置文件
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化 HashMap
	HashMap := uhash.NewHash()
	err = HashMap.LoadFromFile("config/hash.json")
	if err != nil {
		global.LogZap.Error("加载 hash 文件失败", err)
		return
	}

	// 初始化 Telegram 配置
	t := Tbot.Telegramini{
		Apitoken:   cfg.Section("telegram").Key("telegram_apitoken").String(),
		TgimageUrl: cfg.Section("telegram").Key("image_farm").String(),
		Notifyurl:  cfg.Section("telegram").Key("notify_url").String(),
		Hash:       HashMap,
	}

	// 创建并启动 Telegram Bot
	bot, err := Tbot.NewTeleBot(&t)
	if err != nil {
		log.Fatalf("启动 bot 失败: %v", err)
	}

	// 设置路由
	r := router.SetupRoutes(HashMap, bot)

	// 启动 Telegram Bot
	go bot.Bot.Start()

	// 日志记录启动信息
	global.LogZap.Infoln("启动成功", t.Apitoken, t.TgimageUrl, t.Notifyurl)

	// 启动 HTTP 服务器
	if err := r.Run(":2095"); err != nil {
		global.LogZap.Fatalf("启动 HTTP 服务器失败: %v", err)
	}
}
