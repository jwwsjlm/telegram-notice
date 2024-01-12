package types

import (
	uhash "telegram-notice/hash"
)

type Result struct {
	File_id        string `json:"file_id"`
	File_unique_id string `json:"file_unique_id"`
	File_size      int    `json:"file_size"`
	File_path      string `json:"file_path"`
}

type Telegramini struct {
	Apitoken   string
	TgimageUrl string
	Hash       *uhash.Hashtable
}
type Data struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}
type Image struct {
	Src string `json:"src"`
}
