package utils

import (
	"github.com/imroc/req/v3"
	"io"
	"log"
	types "telegram-notice/struct"
)

var client = req.Client{}

func init() {
	client = *req.C()
}
func Upimage(u string, data []byte) (string, error) {
	var imgt []types.Image

	_, err := client.R().
		SetFileBytes("file", "test.jpeg", data).
		SetHeaders(map[string]string{
			"authority": u,
			"referer":   "https://" + "+url.Host+" + "/",
			"accept":    "application/json, text/plain, */*",
		}).
		SetSuccessResult(&imgt).
		Post("https://" + u + "/upload")
	if err != nil {

		return "", err
	}
	ret := "https://" + u + imgt[0].Src
	log.Println("来新图啦:", ret)
	return ret, nil
}
func GetImage(u string) ([]byte, error) {
	resp, err := client.R().Get(u)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error when closing the body: %v", err)
		}
	}(resp.Body)
	return resp.Bytes(), nil
}
func GetFile(u string) (types.Data, error) {
	var ret types.Data
	_, err := client.R().
		SetSuccessResult(&ret).
		Get(u)
	if err != nil {
		return types.Data{}, err
	}
	return ret, nil
}
