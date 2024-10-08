package Tbot

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dustin/go-humanize"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"telegram-notice/global"
	uhash "telegram-notice/hash"
	"telegram-notice/utils"
)

// Telegramini 结构体用于存储 Telegram 相关配置和 Bot 实例
type Telegramini struct {
	Apitoken   string
	TgimageUrl string
	Notifyurl  string
	Hash       *uhash.Hashtable
	Bot        *tele.Bot
}

// commandGethook 处理 /gethook 命令，返回用户的 MD5 ID
func (tm *Telegramini) commandGethook(tc tele.Context) error {
	md5id := uhash.IntToMd5(tc.Message().Chat.ID)
	txt := fmt.Sprintf("您的用户ID为\n```%s```\n请复制到网页中使用\nhttps://%s/webhook/%s?text=hello\n请勿泄露您的md5给他人.否则可能导致他人发送消息到您的账号\n未经md5加密之前的id为\n```%d```",
		md5id, tm.Notifyurl, md5id, tc.Message().Chat.ID)
	log.Println("来新用户啦:", tc.Message().Chat.ID, "\n", "md5:", md5id)

	options := &tele.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tele.ModeMarkdown,
	}
	if err := tc.Send(txt, options); err != nil {
		global.LogZap.Errorf("消息发送失败: %v", err)
		return err
	}

	tm.Hash.Set(tc.Message().Chat.ID, md5id)
	if err := tm.Hash.SaveToFile("config/hash.json"); err != nil {
		global.LogZap.Errorf("hash保存失败: %v", err)
		return err
	}
	return nil
}

// procPhoto 处理接收到的照片
func (tm *Telegramini) procPhoto(tc tele.Context) error {
	if tc.Message().Photo.FileSize > 5242880 {
		return tc.Send("文件过大,请上传小于5M的图片或视频")
	}

	tc.Send("正在下载图片,文件大小为:" + humanize.IBytes(uint64(tc.Message().Photo.FileSize)))
	reader, err := tm.Bot.File(&tc.Message().Photo.File)
	if err != nil {
		global.LogZap.Errorf("下载图片失败: %v", err)
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		global.LogZap.Errorf("读取图片失败: %v", err)
		return err
	}

	upimage, err := utils.Upimage(tm.TgimageUrl, data)
	if err != nil {
		global.LogZap.Errorf("上传图片失败: %v", err)
		tc.Send("上传图片失败:" + err.Error())
		return err
	}

	tc.Send("您的图片已经上传到图床\n" + upimage)
	global.LogZap.Infof("来新图啦: %s", upimage)
	return nil
}

// procVideo 处理接收到的视频
func (tm *Telegramini) procVideo(tc tele.Context) error {
	if tc.Message().Video.FileSize > 5242880 {
		return tc.Send("文件过大,请上传小于5M的图片或视频")
	}

	tc.Send("正在下载图片或视频,文件大小为:" + humanize.IBytes(uint64(tc.Message().Video.FileSize)))
	reader, err := tm.Bot.File(&tc.Message().Video.File)
	if err != nil {
		global.LogZap.Errorf("下载视频失败: %v", err)
		tc.Send("下载图片失败:" + err.Error())
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		global.LogZap.Errorf("读取视频失败: %v", err)
		return err
	}

	upimage, err := utils.Upimage(tm.TgimageUrl, data)
	if err != nil {
		global.LogZap.Errorf("上传视频失败: %v", err)
		tc.Send("上传图片失败:" + err.Error())
		return err
	}

	tc.Send("您的图片已经上传到图床\n" + upimage)
	global.LogZap.Infof("来新图啦: %s", upimage)
	return nil
}

// procAnimation 处理接收到的动画
func (tm *Telegramini) procAnimation(tc tele.Context) error {
	if tc.Message().Animation.FileSize > 5242880 {
		return tc.Send("文件过大,请上传小于5M的图片或视频")
	}

	tc.Send("正在下载图片,文件大小为:" + humanize.IBytes(uint64(tc.Message().Animation.FileSize)))
	reader, err := tm.Bot.File(&tc.Message().Animation.File)
	if err != nil {
		global.LogZap.Errorf("下载动画失败: %v", err)
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		global.LogZap.Errorf("读取动画失败: %v", err)
		return err
	}

	upimage, err := utils.Upimage(tm.TgimageUrl, data)
	if err != nil {
		global.LogZap.Errorf("上传动画失败: %v", err)
		tc.Send("上传图片失败:" + err.Error())
		return err
	}

	tc.Send("您的图片已经上传到图床\n" + upimage)
	global.LogZap.Infof("来新图啦: %s", upimage)
	return nil
}

// procDocument 处理接收到的文档
func (tm *Telegramini) procDocument(tc tele.Context) error {
	if tc.Message().Document.MIME != "image/jpeg" && tc.Message().Document.MIME != "image/png" {
		return tc.Send("不支持的文件类型")
	} else if tc.Message().Document.FileSize > 5242880 {
		return tc.Send("文件过大,请上传小于5M的图片或视频")
	}

	err := tc.Send("正在下载图片,文件大小为:" + humanize.IBytes(uint64(tc.Message().Document.FileSize)))
	if err != nil {
		return err
	}
	reader, err := tm.Bot.File(&tc.Message().Document.File)
	if err != nil {
		global.LogZap.Errorf("下载文档失败: %v", err)
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		global.LogZap.Errorf("读取文档失败: %v", err)
		return err
	}

	upimage, err := utils.Upimage(tm.TgimageUrl, data)
	if err != nil {
		global.LogZap.Errorf("上传文档失败: %v", err)
		tc.Send("上传图片失败:" + err.Error())
		return err
	}

	err = tc.Send("您的图片已经上传到图床\n" + upimage)
	if err != nil {
		return err
	}
	global.LogZap.Infof("来新图啦: %s", upimage)
	return nil
}

// procOnText 处理接收到的文本消息
func (tm *Telegramini) procOnText(tc tele.Context) error {
	return tc.Send("您的消息已经收到\n" + tc.Message().Text)
}

// setRouters 设置 Telegram Bot 的命令和处理器
func (tm *Telegramini) setRouters() {
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

	tm.Bot.Use(middleware.AutoRespond())
	tm.Bot.Handle("/gethook", tm.commandGethook, middleware.IgnoreVia())
	tm.Bot.Handle("/help", func(c tele.Context) error {
		return c.Send("输入/gethook获得您的用户ID\n" + "悄悄告诉你.直接发送图片也可以上传哦")
	}, middleware.IgnoreVia())
	tm.Bot.Handle("/start", func(c tele.Context) error {
		return c.Send("悄悄告诉你,直接发送图片给我,可以直接上传图片到telegram图床\n"+"本项目开源地址:https://github.com/jwwsjlm/telegram-notice", menu)
	}, middleware.IgnoreVia())
	tm.Bot.Handle(tele.OnDocument, tm.procDocument, middleware.IgnoreVia())
	tm.Bot.Handle(tele.OnPhoto, tm.procPhoto, middleware.IgnoreVia())
	tm.Bot.Handle(tele.OnVideo, tm.procVideo, middleware.IgnoreVia())
	tm.Bot.Handle(tele.OnText, tm.procOnText, middleware.IgnoreVia())
	tm.Bot.Handle(tele.OnAnimation, tm.procAnimation, middleware.IgnoreVia())
}

// NewTeleBot 创建并初始化一个新的 Telegram Bot
func NewTeleBot(t *Telegramini) (*Telegramini, error) {
	pref := tele.Settings{
		Token:  t.Apitoken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		global.LogZap.Fatal(err)
		return nil, err
	}

	t.Bot = b
	t.setRouters()
	return t, nil
}
