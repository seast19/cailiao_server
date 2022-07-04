package models

import (
	"cailiao_server/utils"
	"errors"
	"fmt"
	"regexp"
)

const salt = "qwer12345"

// UserAddUser 添加用户
func UserAddUser(newUser *User) (uint, error) {
	//检查参数
	matched, err := regexp.Match(`\d{11}`, []byte(newUser.Phone))
	if !matched {
		return 0, errors.New("用户手机号错误")
	}
	if !utils.ValidateRole(newUser.Role) {
		return 0, errors.New("用户权限设置错误")
	}

	// 填充信息
	newUser.Password = utils.Md5(newUser.Password + salt) //加密密码

	err = GlobalDb.Create(&newUser).Error
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("用户已存在")
	}

	fmt.Printf("创建用户成功 [%d]%s\n", newUser.ID, newUser.Phone)

	return newUser.ID, nil

}

// UserDelUserById 删除用户
func UserDelUserById(id uint) (bool, error) {

	user := User{}
	user.ID = id

	err := GlobalDb.Delete(&user).Error
	if err != nil {
		fmt.Println(err)
		return false, errors.New("删除用户失败")
	}

	return true, nil

}

// UserGetUserByPhone 根据phone获取用户信息
func UserGetUserByPhone(phone string) (*User, error) {

	user := User{}
	err := GlobalDb.Preload("Car").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("查询数据库失败")
	}

	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

// UserGetUserById 根据id获取用户信息
func UserGetUserById(id uint) (*User, error) {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return nil, err
	//}
	//defer GLOBAL_DB.Close()

	user := User{}
	err := GlobalDb.Preload("Car").Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("查询数据库失败")
	}

	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

//检验密码是否相同
func UserIsPwdSame(pwd, pwdMD5 string) bool {
	if utils.Md5(pwd+salt) == pwdMD5 {
		return true
	}
	return false
}

// UserCheckPwd 检查密码正确
func UserCheckPwd(phone, pwd string) (bool, error) {

	//	检查密码是否正确
	pwdMd5 := utils.Md5(pwd + salt)
	//fmt.Println(phone,pwd)
	user := User{}
	var count int64 = 0
	err := GlobalDb.Model(&user).Where("phone = ? AND password = ?", phone, pwdMd5).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return false, errors.New("查询数据库失败")
	}

	if count == 1 {
		return true, nil
	}

	return false, nil
}

// UserGetUsersByPage 分页获取用户
func UserGetUsersByPage(page, perPage int) ([]User, int64, error) {

	var users []User
	var count int64 = 0

	err := GlobalDb.Preload("Car").Omit("password").Offset(perPage * (page - 1)).Limit(perPage).Order("id DESC").Find(&users).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询数据库失败")
	}
	err = GlobalDb.Model(&User{}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询数据库失败")
	}

	return users, count, nil
}

// UserUpdateUserById 修改用户
func UserUpdateUserById(u *User) error {

	//获取原始用户
	user, err := UserGetUserById(u.ID)
	if err != nil {
		return err
	}

	//检查参数，防止直接插入有问题数据
	if len(u.Phone) > 0 {
		matched, _ := regexp.Match(`\d{11}`, []byte(u.Phone))
		if !matched {
			return errors.New("用户手机号格式错误")
		}

		user.Phone = u.Phone
	}

	if len(u.RealName) > 0 {
		user.RealName = u.RealName
	}

	if len(u.Password) > 0 {
		user.Password = utils.Md5(u.Password + salt) //加密密码
	}

	if len(u.Role) > 0 {
		if !utils.ValidateRole(u.Role) {
			return errors.New("用户权限设置错误")
		}

		user.Role = u.Role
	}
	user.CarID = u.CarID

	// 填充信息

	err = GlobalDb.Model(u).Updates(&user).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("更新数据库失败")
	}

	return nil
}

// UserGetIDByJwt 根据jwt获取用户id
func UserGetIDByJwt(jwt string) (uint, error) {
	phone, _, err := utils.ParseJWT(jwt)
	if err != nil {
		return 0, err
	}

	user, err := UserGetUserByPhone(phone)
	return user.ID, err

}
