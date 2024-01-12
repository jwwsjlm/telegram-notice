package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"log"
	"telegram-notice/hash"
	"telegram-notice/router"
	types "telegram-notice/struct"
	"telegram-notice/tgbot"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		return
	}
	var Tini types.Telegramini
	HashMap := uhash.NewHash()
	err = HashMap.LoadFromFile("config/hash.json")
	Tini.Apitoken = cfg.Section("telegram").Key("telegram_apitoken").String()
	Tini.TgimageUrl = cfg.Section("telegram").Key("image_farm").String()
	Tini.Hash = HashMap

	if err != nil {
		log.Panicln("加载hash文件失败", err)
		return
	}
	bot := tgbot.NewBot(Tini)
	bot.Hash = HashMap
	r := router.SetupRoutes(HashMap, bot)
	err = r.Run(":2095")
	if err != nil {
		return
	}
}
