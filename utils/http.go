package utils

import (
	"github.com/imroc/req/v3"
	"io"
	"log"
	types "telegram-notice/struct"
)

// 初始化 HTTP 客户端
var client = req.Client{}

func init() {
	client = *req.C()
}

// Upimage 上传图片到指定 URL
func Upimage(u string, data []byte) (string, error) {
	var imgt []types.Image

	// 发起 POST 请求上传图片
	_, err := client.R().
		SetFileBytes("file", "test.jpeg", data). // 设置文件字段
		SetHeaders(map[string]string{
			"authority": u,
			"referer":   "https://" + u + "/",
			"accept":    "application/json, text/plain, */*",
		}).
		SetSuccessResult(&imgt). // 设置成功响应的结果类型
		Post("https://" + u + "/upload")
	if err != nil {
		return "", err
	}

	// 返回上传后的图片 URL
	ret := "https://" + u + imgt[0].Src
	log.Println("来新图啦:", ret)
	return ret, nil
}

// GetImage 从指定 URL 获取图片数据
func GetImage(u string) ([]byte, error) {
	resp, err := client.R().Get(u)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("关闭响应体时出错: %v", err)
		}
	}(resp.Body)

	// 返回图片数据
	return resp.Bytes(), nil
}

// GetFile 从指定 URL 获取文件数据
func GetFile(u string) (types.Data, error) {
	var ret types.Data

	// 发起 GET 请求获取文件数据
	_, err := client.R().
		SetSuccessResult(&ret). // 设置成功响应的结果类型
		Get(u)
	if err != nil {
		return types.Data{}, err
	}

	// 返回文件数据
	return ret, nil
}
