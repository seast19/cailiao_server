package models

import (
	"cailiao_server/utils"
	"errors"
	"fmt"
)

const salt = "qwer12345"

// AddUser 添加用户
// return
func AddUser(newUser *User) (uint, error) {
	db, err := getConn()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx := db.Begin()

	// 检查用户是否存在，不存在则可以添加用户
	user := User{}
	tx.Where(&User{Phone: newUser.Phone}).First(&user)
	if user.Phone != "" {
		tx.Rollback()
		fmt.Println("用户已存在")
		return 0, errors.New("用户已存在")
	}

	// 填充信息
	newUser.Password = utils.Md5(newUser.Phone + newUser.Password + salt) //加密密码

	tx.Create(newUser)

	tx.Commit()
	fmt.Printf("创建用户成功 [%d]%s\n", newUser.ID, newUser.Phone)

	return newUser.ID, nil

}

// 删除用户
// 修改用户

// 用户登录
// 校验jwt
