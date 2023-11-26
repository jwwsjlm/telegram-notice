package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"telegram-notice/hash"
	"telegram-notice/router"
	"telegram-notice/tgbot"
)

var HashMap *uhash.Hashtable

type Config struct {
	TELEGRAM_APITOKEN string
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		return
	}
	TELEGRAM_APITOKEN := cfg.Section("telegram").Key("TELEGRAM_APITOKEN").String()
	fmt.Println("Telegram API Token:", TELEGRAM_APITOKEN)

	HashMap = uhash.Newhash()
	err = HashMap.LoadFromFile("./hash.json")
	if err != nil {
		return
	}
	bot := tgbot.NewBot(TELEGRAM_APITOKEN, HashMap)

	//使用gin创建一个路由 用于接收telegram的webhook
	webhook := gin.Default()

	router.SetupRoutes(webhook, HashMap, *bot)

	err = webhook.Run(":8080")
	if err != nil {
		return
	}

}
