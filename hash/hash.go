package uhash

import (
	"github.com/goccy/go-json"
	"os"

	"telegram-notice/utils"
)

type Hashtable struct {
	hashMap map[string]int64
}

func (h Hashtable) Get(string2 string) (int64, bool) {
	if _, ok := h.hashMap[string2]; !ok {
		return 0, false
	}
	return h.hashMap[string2], true
}
func (h Hashtable) Set(int642 int64) {
	h.hashMap[utils.IntToMd5(int642)] = int642
}
func (h Hashtable) SaveToFile(filename string) error {
	jsonData, err := json.Marshal(h.hashMap)
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
func (h Hashtable) LoadFromFile(filename string) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 将 JSON 数据解析为哈希表
	var data map[string]int64
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}
	h.hashMap = data
	return nil
}
func Newhash() *Hashtable {
	return &Hashtable{hashMap: make(map[string]int64)}
}
