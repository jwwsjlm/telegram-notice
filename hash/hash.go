package uhash

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"strconv"
	"sync"
)

type Hashtable struct {
	hashMap *sync.Map
}

func IntToMd5(i int64) string {
	str := strconv.FormatInt(i, 10)
	hash := md5.Sum([]byte(str))
	md5Str := fmt.Sprintf("%x", hash)
	return md5Str
}

func (h *Hashtable) Get(string2 string) (int64, error) {
	val, ok := h.hashMap.Load(string2)
	if !ok {
		return 0, errors.New("未找到当前用户id")
	}
	result, ok := val.(int64)
	if !ok {
		return 0, errors.New("转换失败联系管理员")
	} else {
		return result, nil
	}

}
func (h *Hashtable) Set(int642 int64, md5 string) {

	h.hashMap.Store(md5, int642)

}

func (h *Hashtable) SaveToFile(filename string) error {
	// 将 sync.Map 转换为标准的 Go 数据结构
	data := make(map[string]interface{})
	h.hashMap.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})

	// 将数据转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 将 JSON 数据写入文件
	err = os.WriteFile(filename, jsonData, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hashtable) LoadFromFile(filename string) error {

	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	if fileInfo.Size() == 0 {
		// 文件内容为空，不进行读取操作
		return nil
	}
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 将 JSON 数据解析为标准的 Go 数据结构
	var data map[string]int64
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	// 将数据存储到 sync.Map 中
	h.hashMap.Range(func(key, _ interface{}) bool {
		h.hashMap.Delete(key)
		return true
	})
	for key, value := range data {
		h.hashMap.Store(key, value)
	}

	return nil
}
func NewHash() *Hashtable {

	return &Hashtable{hashMap: &sync.Map{}}
}
