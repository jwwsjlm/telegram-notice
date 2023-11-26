package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"telegram-notice/hash"
	"telegram-notice/router"
	"telegram-notice/tgbot"
)

var HashMap *uhash.Hashtable

func main() {
	HashMap = uhash.Newhash()
	err := HashMap.LoadFromFile("hash.json")
	if err != nil {
		return
	}
	bot := tgbot.NewBot(os.Getenv("TELEGRAM_APITOKEN"), HashMap)

	//使用gin创建一个路由 用于接收telegram的webhook
	webhook := gin.Default()

	router.SetupRoutes(webhook, HashMap, *bot)

	err = webhook.Run(":8080")
	if err != nil {
		return
	}

}
