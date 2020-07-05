package utils

import (
	"crypto/md5"
	"fmt"
)

// Md5 生成MD5
func Md5(s string) string {
	data := []byte(s)
	r := fmt.Sprintf("%x", md5.Sum(data))
	return r
}
