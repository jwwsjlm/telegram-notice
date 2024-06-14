package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"log"
	"os"
	Tbot "telegram-notice/Newtele"
	"telegram-notice/hash"
	"telegram-notice/router"
)

func main() {
	gin.SetMode(gin.TestMode)
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		panic(err)
		return
	}
	logFile, err := os.OpenFile("config/telebot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)
	logger := log.New(logFile, "telebot: ", log.LstdFlags)
	HashMap := uhash.NewHash()
	err = HashMap.LoadFromFile("config/hash.json")
	t := Tbot.Telegramini{
		Apitoken:   cfg.Section("telegram").Key("telegram_apitoken").String(),
		TgimageUrl: cfg.Section("telegram").Key("image_farm").String(),
		Notifyurl:  cfg.Section("telegram").Key("notify_url").String(),
		Hash:       HashMap,
		Log:        logger,
	}

	if err != nil {
		log.Panicln("加载hash文件失败", err)
		return
	}

	bot, err := Tbot.NewTeleBot(&t)
	if err != nil {
		log.Panicln("启动bot失败", err)
		return
	}

	r := router.SetupRoutes(HashMap, bot)
	go bot.Bot.Start()
	log.Println("启动成功", "\n", t.Apitoken, "\n", t.TgimageUrl, "\n", t.Notifyurl)
	err = r.Run(":2095")

	if err != nil {
		logger.Println("启动失败", err)
		return
	}

}
