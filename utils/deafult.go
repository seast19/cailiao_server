package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

// Md5 使用string生成MD5
func Md5(s string) string {
	data := []byte(s)
	r := fmt.Sprintf("%x", md5.Sum(data))
	return r
}

// IndexOfStringList 判断某元素是否在[]string数组中
func IndexOfStringList(arr []string, target string) int {
	for index, value := range arr {
		if value == target {
			return index
		}
	}
	return -1
}

// ValidateRole 校验用户角色字符串是否合法，合法参数[user editor admin]
func ValidateRole(s string) bool {
	roles := []string{
		"user",
		"editor",
		"admin",
	}
	for _, value := range roles {
		if value == s {
			return true
		}
	}
	return false
}

// StringToInt 将字符转为数字int，默认为0
func StringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
