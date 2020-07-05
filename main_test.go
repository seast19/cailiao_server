package main

import (
	"cailiao_server/models"
	"testing"
)

// func TestJwt(t *testing.T) {
// 	utils.GenJWT()
// }

// func TestParseJWT(t *testing.T) {
// 	utils.ParseJWT("sds")
// }

func TestAddUser(t *testing.T) {
	user := models.User{
		Phone:    "12323ew2322",
		Password: "qwer",
		Role:     "admin",
	}

	models.AddUser(&user)
}
