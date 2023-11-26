package uhash

import "telegram-notice/utils"

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
func Newhash() *Hashtable {
	return &Hashtable{hashMap: make(map[string]int64)}
}
