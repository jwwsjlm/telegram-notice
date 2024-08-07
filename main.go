package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"log"
	Tbot "telegram-notice/Newtele"
	"telegram-notice/global"
	"telegram-notice/hash"
	"telegram-notice/router"
)

func main() {
	gin.SetMode(gin.TestMode)
	InitLogger()
	// 程序结束前同步日志缓冲区
	defer global.LogZap.Sync()
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		panic(err)
		return
	}

	HashMap := uhash.NewHash()
	err = HashMap.LoadFromFile("config/hash.json")
	t := Tbot.Telegramini{
		Apitoken:   cfg.Section("telegram").Key("telegram_apitoken").String(),
		TgimageUrl: cfg.Section("telegram").Key("image_farm").String(),
		Notifyurl:  cfg.Section("telegram").Key("notify_url").String(),
		Hash:       HashMap,
	}

	if err != nil {

		global.LogZap.Error("加载hash文件失败", err)
		return
	}

	bot, err := Tbot.NewTeleBot(&t)
	if err != nil {
		log.Panicln("启动bot失败", err)
		return
	}

	r := router.SetupRoutes(HashMap, bot)
	go bot.Bot.Start()
	global.LogZap.Infoln("启动成功", "\n", t.Apitoken, "\n", t.TgimageUrl, "\n", t.Notifyurl)

	err = r.Run(":2095")

	if err != nil {
		global.LogZap.Infoln("启动失败", err)
		return
	}

}
