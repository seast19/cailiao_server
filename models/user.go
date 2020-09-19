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
	db, err := getConn()
	if err != nil {
		return 0, err
	}
	defer db.Close()

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

	err = db.Create(&newUser).Error
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("用户已存在")
	}

	fmt.Printf("创建用户成功 [%d]%s\n", newUser.ID, newUser.Phone)

	return newUser.ID, nil

}

// 删除用户
func UserDelUserById(id uint) (bool, error) {
	db, err := getConn()
	if err != nil {
		return false, err
	}
	defer db.Close()

	user := User{}
	user.ID = id

	err = db.Delete(&user).Error
	if err != nil {
		fmt.Println(err)
		return false, errors.New("删除用户失败")
	}

	return true, nil

}

//根据phone获取用户信息
func UserGetUserByPhone(phone string) (*User, error) {
	db, err := getConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	user := User{}
	err = db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("查询数据库失败")
	}

	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

//根据id获取用户信息
func UserGetUserById(id uint) (*User, error) {
	db, err := getConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	user := User{}
	err = db.Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("查询数据库失败")
	}

	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

// 检查密码正确
func UserCheckPwd(phone, pwd string) (bool, error) {
	db, err := getConn()
	if err != nil {
		return false, err
	}
	defer db.Close()

	//	检查密码是否正确
	pwdMd5 := utils.Md5(pwd + salt)

	user := User{}
	count := 0
	err = db.Model(&user).Where("phone = ? AND password >= ?", phone, pwdMd5).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return false, errors.New("查询数据库失败")
	}

	if count == 1 {
		return true, nil
	}

	return false, nil
}

//分页获取用户
func UserGetUsersByPage(page, perPage int) ([]User, int, error) {
	db, err := getConn()
	if err != nil {
		return nil, 0, err
	}
	defer db.Close()

	var users []User
	count := 0

	err = db.Select("phone, nick_name, real_name,role,ID").Offset(perPage * (page - 1)).Limit(perPage).Order("id DESC").Find(&users).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询数据库失败")
	}
	err = db.Model(&User{}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询数据库失败")
	}

	return users, count, nil
}

//修改用户
func UserUpdateUserById(u *User) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	//获取原始用户
	user, err := UserGetUserById(u.ID)
	if err != nil {
		return err
	}

	//检查参数
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
	// 填充信息

	err = db.Model(u).Update(&user).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("更新数据库失败")
	}

	return nil
}

//构建jwt获取用户id
func UserGetIDByJwt(jwt string) (uint,error) {
	phone, err := utils.ParseJWT(jwt)
	if err!=nil{
		return 0,err
	}

	user,err:= UserGetUserByPhone(phone)
	return user.ID,err

}
