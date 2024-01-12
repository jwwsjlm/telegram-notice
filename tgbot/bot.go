package tgbot

import (
	"fmt"
	"github.com/dustin/go-humanize"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	uhash "telegram-notice/hash"
	types "telegram-notice/struct"
	"telegram-notice/utils"
)

type TgBot struct {
	Bot     *tgbotapi.BotAPI
	u       tgbotapi.UpdateConfig
	UimgUrl types.Telegramini
	Hash    *uhash.Hashtable
}

func SendMsg(t *TgBot, c int64, text string) error {
	msg := tgbotapi.NewMessage(c, text)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
func MessagePhoto(t *TgBot, update *tgbotapi.Update) error {
	err := SendMsg(t, update.Message.Chat.ID, "正在上传图片,文件大小为:"+humanize.IBytes(uint64(update.Message.Photo[3].FileSize)))
	if err != nil {
		return err
	}
	if update.Message.Photo[3].FileSize > 0 {
		fileId := update.Message.Photo[3].FileID
		ret, err := utils.GetFile("https://api.telegram.org/bot" + t.Bot.Token + "/getFile?file_id=" + fileId)
		if err != nil {
			errs := fmt.Errorf("消息发送失败: %w", err)
			return errs

		}
		img, err := utils.GetImage("https://api.telegram.org/file/bot" + t.Bot.Token + "/" + ret.Result.File_path)
		if err != nil {
			errs := fmt.Errorf("GetImage失败: %w", err)
			return errs

		}
		r, err := utils.Upimage(t.UimgUrl.TgimageUrl, img)

		if err != nil {
			errs := fmt.Errorf("Upimage失败: %w", err)
			return errs
		}
		err = SendMsg(t, update.Message.Chat.ID, "您的图片已经上传到图床\n"+r)
		log.Println("来新图啦:", r)
		if err != nil {
			errs := fmt.Errorf("发送消息失败: %w", err)
			return errs
		}
		//
	}
	return nil
}
func MessageDocument(t *TgBot, update *tgbotapi.Update) error {
	err := SendMsg(t, update.Message.Chat.ID, "正在上传图片,文件大小为:"+humanize.IBytes(uint64(update.Message.Document.FileSize)))
	if err != nil {
		errs := fmt.Errorf("发送消息失败: %w", err)
		return errs
	}
	if update.Message.Document.FileSize > 0 {
		fileId := update.Message.Document.FileID

		ret, err := utils.GetFile("https://api.telegram.org/bot" + t.Bot.Token + "/getFile?file_id=" + fileId)
		if err != nil {
			errs := fmt.Errorf("获取File信息失败: %w", err)
			return errs

		}
		img, err := utils.GetImage("https://api.telegram.org/file/bot" + t.Bot.Token + "/" + ret.Result.File_path)
		if err != nil {
			errs := fmt.Errorf("获取文件信息失败: %w", err)
			return errs
		}
		r, err := utils.Upimage(t.UimgUrl.TgimageUrl, img)

		if err != nil {
			errs := fmt.Errorf("上传图片失败失败: %w", err)
			return errs
		}
		err = SendMsg(t, update.Message.Chat.ID, "您的图片已经上传到图床\n"+r)
		log.Println("来新图啦:", r)
		if err != nil {
			errs := fmt.Errorf("消息发送失败: %w", err)
			return errs

		}
	}
	return nil
}
func CommandGethook(t *TgBot, c int64) error {
	md5id := uhash.IntToMd5(c)
	log.Println("md5id", md5id, "update.Message.Chat.ID", c)
	err := SendMsg(t, c, "您的用户ID为\n"+md5id)
	if err != nil {
		errs := fmt.Errorf("消息发送失败: %w", err)
		return errs

	}
	//update.Message.Chat.ID转换为string类型
	t.Hash.Set(c, md5id)
	err = t.Hash.SaveToFile("config/hash.json")
	if err != nil {
		errs := fmt.Errorf("hash保存失败: %w", err)
		return errs

	}
	return nil
}
func MsgProc(t *TgBot, update *tgbotapi.Update) {

	if update.Message != nil {
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "选择你的功能:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("help", "/help"),
					tgbotapi.NewInlineKeyboardButtonData("gethook", "/gethook"),
				),
			)
			t.Bot.Send(msg)
		}
		if update.Message.IsCommand() { // ignore any non-command Messages
			switch update.Message.Command() {
			case "gethook":
				err := CommandGethook(t, update.Message.Chat.ID)
				if err != nil {
					errs := SendMsg(t, update.Message.Chat.ID, "失败:"+err.Error())
					if errs != nil {
						log.Println("消息发送失败:", errs.Error())
					}
				}
				return
			default:
				err := SendMsg(t, update.Message.Chat.ID, "输入/gethook获得您的用户ID\n"+"悄悄告诉你.直接发送图片也可以上传哦")
				if err != nil {
					log.Println("消息发送失败gethook:", err.Error())
				}
				return
			}
		}
	}

	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data == "/help" {
			err := SendMsg(t, update.CallbackQuery.Message.Chat.ID, "输入/gethook获得您的用户ID\n"+"悄悄告诉你.直接发送图片也可以上传哦")
			if err != nil {
				log.Println("消息发送失败gethook:", err.Error())
			}
			return
		}
		if update.CallbackQuery.Data == "/gethook" {
			err := CommandGethook(t, update.CallbackQuery.Message.Chat.ID)
			if err != nil {
				errs := SendMsg(t, update.CallbackQuery.Message.Chat.ID, "失败:"+err.Error())
				if errs != nil {
					log.Println("消息发送失败:", errs.Error())
				}

			}
			return
		}

		log.Println("update.CallbackQuery", update.CallbackQuery.Data)
		return
	}

	if update.Message.Photo != nil {
		err := MessagePhoto(t, update)
		if err != nil {
			errs := SendMsg(t, update.Message.Chat.ID, "失败:"+err.Error())
			if errs != nil {
				log.Println("消息发送失败:", errs.Error())
			}

		}
		return
	}
	if update.Message.Document != nil {
		err := MessageDocument(t, update)
		if err != nil {
			errs := SendMsg(t, update.Message.Chat.ID, "失败:"+err.Error())
			if errs != nil {
				log.Println("消息发送失败:", errs.Error())
			}

		}
		return
	}

}

func FormatMessage(t *TgBot) {
	updates := t.Bot.GetUpdatesChan(t.u)
	for update := range updates {
		go MsgProc(t, &update)

	}
}

func (t *TgBot) SendMees(id int64, string2 string) error {
	msg := tgbotapi.NewMessage(id, string2)
	msg.Text = string2
	_, err := t.Bot.Send(msg)
	return err

}
func NewBot(t types.Telegramini) *TgBot {
	bot, err := tgbotapi.NewBotAPI(t.Apitoken)
	log.Println("bot", t)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	tgBot := TgBot{
		Bot:     bot,
		u:       u,
		UimgUrl: t,
		Hash:    t.Hash,
	}
	go FormatMessage(&tgBot)
	return &tgBot
}
