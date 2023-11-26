package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func IntToMd5(i int64) string {
	str := strconv.FormatInt(i, 10)
	hash := md5.Sum([]byte(str))
	md5Str := fmt.Sprintf("%x", hash)
	return md5Str
}
