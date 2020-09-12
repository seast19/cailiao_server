package utils

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const secretKey = "wqdqwdqqwsqs"

// GenJWT 生成jwt
func GenJWT(phone string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
		"time":  time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secretKey))

	//fmt.Println(tokenString, err)
	return tokenString, err
}

// ParseJWT 解析jwt
func ParseJWT(s string) (string, error) {

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//校验时间
		thatTime := claims["time"].(float64)

		//fmt.Printf("%#v", int64(thatTime))

		if time.Now().Unix()-int64(thatTime) > 15*24*60*60 {
			fmt.Println("令牌过期")
			return "", errors.New("令牌已过期")
		}
		//fmt.Println(claims["foo"], claims["nbf"])
		return claims["phone"].(string), nil
	}
	//fmt.Println(err)
	return "", err

}
