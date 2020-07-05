package utils

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

const secretKey = "wqdqwdqqwsqs"

// GenJWT 生成jwt
func GenJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secretKey))

	fmt.Println(tokenString, err)
	return tokenString, err
}

// ParseJWT 解析jwt
func ParseJWT(s string) (string, error) {
	// sample token string taken from the New example
	// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.vAEOPIhKwANkwed9OSVzEhMuKwZTwr1ocmuqEwfurmY"

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims["username"].(string), nil
	}
	fmt.Println(err)
	return "", err

}
