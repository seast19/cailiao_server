package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

// Md5 生成MD5
func Md5(s string) string {
	data := []byte(s)
	r := fmt.Sprintf("%x", md5.Sum(data))
	return r
}

// IndexOfString 判断某元素是否在数组中
func IndexOfString(arr []string, target string) int {
	for index, value := range arr {
		if value == target {
			return index
		}
	}
	return -1
}

// ValidateRole 校验用户角色是否合法
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

// String2Number 将字符转为数字int，默认为0
func String2Number(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
