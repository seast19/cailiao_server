package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "bad apple"

// GenJWT 根据phone、role参数构建jwt
func GenJWT(phone, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role":  role,
		"phone": phone,
		"time":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

// ParseJWT 解析jwt返回phone、role、error
func ParseJWT(s string) (string, string, error) {
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", "", errors.New("解析jwt失败")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//校验时间
		thatTime := claims["time"].(float64)
		if time.Now().Unix()-int64(thatTime) > 15*24*60*60 {
			return "", "", errors.New("令牌已过期")
		}
		return claims["phone"].(string), claims["role"].(string), nil
	}
	return "", "", errors.New("解析jwt错误")
}
