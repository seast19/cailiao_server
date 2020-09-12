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

//判断某元素是否在数组中
func IndexOfString(arr []string, target string) int {
	for index, value := range arr {
		if value == target{
			return index
		}
	}
	return -1
}

//校验用户角色是否合法
func ValidateRole(s string)  bool{
	roles := []string{
		"user",
		"editor",
		"admin",
	}

	for _, value := range roles {
		if value == s{
			return true
		}
	}
	return false
}
