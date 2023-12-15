package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"telegram-notice/hash"
	"telegram-notice/router"
	"telegram-notice/tgbot"
)

var HashMap *uhash.Hashtable

func main() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		return
	}
	TelegramApitoken := cfg.Section("telegram").Key("TELEGRAM_APITOKEN").String()
	fmt.Println("Telegram API Token:", TelegramApitoken)
	HashMap = uhash.Newhash()
	err = HashMap.LoadFromFile("config/hash.json")
	if err != nil {
		return
	}
	bot := tgbot.NewBot(TelegramApitoken, HashMap)
	r := router.SetupRoutes(HashMap, bot)
	err = r.Run(":2095")
	if err != nil {
		return
	}
}
