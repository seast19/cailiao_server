package main

import (
	"cailiao_server/models"
	"cailiao_server/utils"
	"fmt"
	"testing"
)

func TestJwt(t *testing.T) {
	utils.GenJWT("sdds")
}

// func TestParseJWT(t *testing.T) {
// 	utils.ParseJWT("sds")
// }

func TestAddUser(t *testing.T) {
	user := models.User{
		Phone:    "user",
		Password: "123456",
		Role:     "user",
	}

	models.UserAddUser(&user)
}

func TestDelUser(t *testing.T) {
	models.UserDelUserById(6)
}

func TestGetUsers(t *testing.T) {
	models.UserGetUsersByPage(2, 10)
}

func TestIndexOf(t *testing.T) {
	roles := []string{
		"a",
		"b",
		"c",
	}
	aa := utils.IndexOfString(roles, "d")
	fmt.Println(aa)
}

func TestPlaceAdd(t *testing.T) {
	place := models.Place{
		Position: "4好贵",
	}
	err := models.PlaceAdd(&place)
	fmt.Println(err)
}

func TestPlaceAll(t *testing.T) {
	//places:=[]models.Place{}
	places, err := models.PlaceAll()
	fmt.Println(places, err)
}

func TestUtileWrite(t *testing.T) {
	//utils.CheckAndCreateLogFile("./sss.txt")
}
