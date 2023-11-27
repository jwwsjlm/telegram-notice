package uhash

import (
	"github.com/goccy/go-json"
	"os"
	"sync"

	"telegram-notice/utils"
)

type Hashtable struct {
	hashMap sync.Map
}

func (h *Hashtable) Get(string2 string) (int64, bool) {
	val, ok := h.hashMap.Load(string2)
	if !ok {
		return 0, false
	}
	result, ok := val.(int64)
	if !ok {
		return 0, false
	} else {
		return result, true
	}

}
func (h *Hashtable) Set(int642 int64) {
	h.hashMap.Store(utils.IntToMd5(int642), int642)

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
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hashtable) LoadFromFile(filename string) error {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// 文件不存在，不进行读取操作
		return nil
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
func Newhash() *Hashtable {
	return &Hashtable{hashMap: sync.Map{}}
}
