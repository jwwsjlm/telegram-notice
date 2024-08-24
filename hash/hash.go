package uhash

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/goccy/go-json"
)

// Hashtable 结构体用于存储用户 ID 和 MD5 映射关系
type Hashtable struct {
	hashMap *sync.Map
}

// IntToMd5 将 int64 转换为 MD5 字符串
func IntToMd5(i int64) string {
	str := strconv.FormatInt(i, 10)
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}

// Get 根据 MD5 字符串获取对应的 int64 用户 ID
func (h *Hashtable) Get(md5Str string) (int64, error) {
	val, ok := h.hashMap.Load(md5Str)
	if !ok {
		return 0, errors.New("未找到当前用户 ID")
	}
	result, ok := val.(int64)
	if !ok {
		return 0, errors.New("转换失败，请联系管理员")
	}
	return result, nil
}

// Set 将 int64 用户 ID 和 MD5 字符串存储到哈希表中
func (h *Hashtable) Set(userID int64, md5Str string) {
	h.hashMap.Store(md5Str, userID)
}

// SaveToFile 将哈希表保存到文件中
func (h *Hashtable) SaveToFile(filename string) error {
	data := make(map[string]interface{})
	h.hashMap.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// LoadFromFile 从文件中加载哈希表
func (h *Hashtable) LoadFromFile(filename string) error {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// 如果文件不存在，创建一个新的文件
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
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

	var data map[string]int64
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	// 清空现有的哈希表数据
	h.hashMap.Range(func(key, _ interface{}) bool {
		h.hashMap.Delete(key)
		return true
	})

	// 将数据存储到哈希表中
	for key, value := range data {
		h.hashMap.Store(key, value)
	}

	return nil
}

// NewHash 创建一个新的 Hashtable 实例
func NewHash() *Hashtable {
	return &Hashtable{hashMap: &sync.Map{}}
}
