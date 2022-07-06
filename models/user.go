package models

import (
	"cailiao_server/utils"
	"errors"
	"github.com/beego/beego/v2/core/logs"
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

	//数据库操作
	err = GlobalDb.Create(&newUser).Error
	if err != nil {
		logs.Error(err)
		return 0, errors.New("用户已存在")
	}

	logs.Info("创建用户成功 [%d]%s\n", newUser.ID, newUser.Phone)

	return newUser.ID, nil
}

// UserDelById 根据id删除用户
func UserDelById(id uint) (bool, error) {
	user := User{}
	user.ID = id

	err := GlobalDb.Delete(&user).Error
	if err != nil {
		logs.Error(err)
		return false, errors.New("删除用户失败")
	}
	return true, nil
}

// UserUpdateById 修改用户
func UserUpdateById(u *User) error {
	//获取原始用户
	user, err := UserGetById(u.ID)
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
		logs.Error(err)
		return errors.New("修改用户数据失败")
	}
	return nil
}

// UserGetByPhone 根据phone获取用户信息
func UserGetByPhone(phone string) (*User, error) {
	user := User{}
	err := GlobalDb.
		Preload("Car").
		Where("phone = ?", phone).
		First(&user).Error
	if err != nil {
		logs.Error(err)
		return nil, errors.New("获取用户信息失败#1")
	}

	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

// UserGetById 根据id获取用户信息
func UserGetById(id uint) (*User, error) {
	user := User{}
	err := GlobalDb.
		Preload("Car").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		logs.Error(err)
		return nil, errors.New("获取用户信息失败#2")
	}

	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

// UserGetByJwt 根据jwt获取用户id
func UserGetByJwt(jwt string) (*User, error) {
	phone, _, err := utils.ParseJWT(jwt)
	if err != nil {
		return nil, errors.New("获取用户信息失败#3")
	}

	user, err := UserGetByPhone(phone)
	return user, err

}

// UserIsPwdSame 检验密码是否相同
func UserIsPwdSame(pwd, pwdMD5 string) bool {
	if utils.Md5(pwd+salt) == pwdMD5 {
		return true
	}
	return false
}

// UserGetAllByPage 分页获取所有用户
func UserGetAllByPage(page, perPage int) ([]User, int64, error) {
	var users []User
	var count int64

	err := GlobalDb.
		Preload("Car").
		Omit("password").
		Offset(perPage * (page - 1)).Limit(perPage).
		Order("id DESC").
		Find(&users).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("获取用户数据失败")
	}
	//条数
	err = GlobalDb.
		Model(&User{}).
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("获取用户数据失败")
	}

	return users, count, nil
}


