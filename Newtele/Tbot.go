package Tbot

import (
	"fmt"
	"github.com/dustin/go-humanize"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"io"
	"log"
	uhash "telegram-notice/hash"
	"telegram-notice/utils"
	"time"
)

type Telegramini struct {
	Apitoken   string
	TgimageUrl string
	Notifyurl  string
	Hash       *uhash.Hashtable
	Bot        *tele.Bot
	Log        *log.Logger
}

func (tm *Telegramini) commandGethook(tc tele.Context) error {
	md5id := uhash.IntToMd5(tc.Message().Chat.ID)
	var txt = "您的用户ID为\n```" + md5id + "```\n请复制到网页中使用" +
		"\nhttps://" + tm.Notifyurl + "/webhook/" + md5id + "?text=hello" +
		"\n请勿泄露您的md5给他人.否则可能导致他人发送消息到您的账号" +
		"未经md5加密之前的id为\n```" + fmt.Sprintf("%d", tc.Message().Chat.ID) + "```"
	log.Println("来新用户啦:", tc.Message().Chat.ID, "\n", "md5:", md5id)
	options := &tele.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tele.ModeMarkdown,
	}
	err := tc.Send(txt, options)
	if err != nil {
		errs := fmt.Errorf("消息发送失败: %w", err)
		tm.Log.Println(errs)
		return errs

	}
	tm.Hash.Set(tc.Message().Chat.ID, md5id)
	err = tm.Hash.SaveToFile("config/hash.json")
	if err != nil {
		errs := fmt.Errorf("hash保存失败: %w", err)
		tm.Log.Println(errs)
		return errs
	}
	return nil
}
func (tm *Telegramini) photoProc(tc tele.Context) error {

	reader, err := tm.Bot.File(&tc.Message().Photo.File)
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			tm.Log.Printf("Error when closing the body: %v", err)
		}
	}(reader)
	if err != nil {
		tm.Log.Println("Failed to download:", err)
		return err
	}

	data, err := io.ReadAll(reader)
	err = tc.Send("正在上传图片,文件大小为:" + humanize.IBytes(uint64(tc.Message().Photo.FileSize)))
	if err != nil {
		tm.Log.Println("发送信息失败:", err)
		return err
	}
	upimage, err := utils.Upimage(tm.TgimageUrl, data)
	if err != nil {
		tm.Log.Println("上传图片失败:", err)
		_ = tc.Send("上传图片失败:" + err.Error())
		return err
	}
	err = tc.Send("您的图片已经上传到图床\n" + upimage)
	if err != nil {
		tm.Log.Println("发送信息失败:", err)
		return err
	}
	tm.Log.Println("来新图啦:", upimage)

	// 响应信息给用户
	return nil

}
func (tm *Telegramini) documentProc(tc tele.Context) error {
	reader, err := tm.Bot.File(&tc.Message().Document.File)
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			tm.Log.Println("Error when closing the body: %v", err)
		}
	}(reader)
	if err != nil {
		tm.Log.Println("Failed to download:", err)
		return err
	}
	data, err := io.ReadAll(reader)

	err = tc.Send("正在上传图片,文件大小为:" + humanize.IBytes(uint64(tc.Message().Document.FileSize)))
	upimage, err := utils.Upimage(tm.TgimageUrl, data)
	if err != nil {
		tm.Log.Println("上传图片失败:", err)
		_ = tc.Send("上传图片失败:" + err.Error())
		return err
	}
	err = tc.Send("您的图片已经上传到图床\n" + upimage)
	if err != nil {
		tm.Log.Println("发送信息失败:", err)
		return err
	}
	tm.Log.Println("来新图啦:", upimage)

	// 响应信息给用户
	return nil

}
func (tm *Telegramini) onTextProc(tc tele.Context) error {
	return tc.Send("您的消息已经收到\n" + tc.Message().Text)
}
func (tm *Telegramini) SetRouters() {
	menu := &tele.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	btnHelp := menu.Text("/help")
	btnSettings := menu.Text("/gethook")
	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)

	tm.Bot.Use(middleware.Logger(tm.Log))
	tm.Bot.Use(middleware.AutoRespond())
	tm.Bot.Handle("/gethook", tm.commandGethook, middleware.IgnoreVia())
	tm.Bot.Handle("/help", func(c tele.Context) error {
		return c.Send("输入/gethook获得您的用户ID\n" + "悄悄告诉你.直接发送图片也可以上传哦")
	}, middleware.IgnoreVia())
	tm.Bot.Handle("/start", func(c tele.Context) error {
		return c.Send("悄悄告诉你,直接发送图片给我,可以直接上传图片到telegram图床", menu)
	}, middleware.IgnoreVia())
	tm.Bot.Handle(tele.OnDocument, tm.documentProc, middleware.IgnoreVia())
	tm.Bot.Handle(tele.OnPhoto, tm.photoProc, middleware.IgnoreVia())

	tm.Bot.Handle(tele.OnText, tm.onTextProc, middleware.IgnoreVia())
}
func NewTeleBot(t *Telegramini) (*Telegramini, error) {
	pref := tele.Settings{
		Token:  t.Apitoken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		t.Log.Fatal(err)
		return nil, err
	}

	t.Bot = b
	t.SetRouters()
	return t, nil
}