package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uhash "telegram-notice/hash"

	"log"
	"telegram-notice/utils"
)

type TgBot struct {
	Bot *tgbotapi.BotAPI
	u   tgbotapi.UpdateConfig
}

func FormatMessage(t TgBot, h *uhash.Hashtable) {

	updates := t.Bot.GetUpdatesChan(t.u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你好呀")
		fmt.Println("目标id", update.Message.Chat.ID)
		// Extract the command from the Message.
		switch update.Message.Command() {
		case "gethook":
			md5id := utils.IntToMd5(update.Message.Chat.ID)
			msg.Text = md5id
			//update.Message.Chat.ID转换为string类型
			h.Set(update.Message.Chat.ID)
			err := h.SaveToFile("./hash.json")
			if err != nil {
				return
			}

		default:
			msg.Text = "输入/gerhook获得您的用户ID"
		}

		if _, err := t.Bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func (t TgBot) SendMees(id int64, string2 string) error {
	msg := tgbotapi.NewMessage(id, string2)
	msg.Text = string2
	_, err := t.Bot.Send(msg)
	return err

}
func NewBot(t string, hash *uhash.Hashtable) *TgBot {
	bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	tgBot := TgBot{
		Bot: bot,
		u:   u,
	}
	go FormatMessage(tgBot, hash)
	return &tgBot
}
