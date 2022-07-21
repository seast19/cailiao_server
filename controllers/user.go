package controllers

import (
	"cailiao_server/models"
	"cailiao_server/utils"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func Hellow(c *gin.Context) {
	c.JSON(http.StatusOK, "this is default route ")
}

// UserLogin 登录
func UserLogin(c *gin.Context) {
	//获取参数
	data := struct {
		Phone string `json:"phone"`
		Pwd   string `json:"pwd"`
	}{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "用户输入参数错误",
		})
		return
	}

	//查询数据库
	user, err := models.UserGetByPhone(data.Phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}
	//fmt.Println(data.Pwd, user.Password)
	//检验用户密码
	if !models.UserIsPwdSame(data.Pwd, user.Password) {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "账户或密码错误",
		})
		return
	}

	//生成jwt
	jwt, err := utils.GenJWT(data.Phone, user.Role)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5000,
			"msg":  "生成jwt失败",
		})
		return
	}

	//返回数据
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "登录成功",
		"data": gin.H{
			"jwt":   jwt,
			"phone": user.Phone,
			"role":  user.Role,
			"car":   user.CarID,
			"user":  user,
		},
	})
}

// UserUpdateById 修改用户
func UserUpdateById(c *gin.Context) {
	idStr := c.Param("id")

	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	data := struct {
		Phone string `json:"phone"`
		Name  string `json:"name"`
		Pwd   string `json:"pwd"`
		Role  string `json:"role"`
		Car   uint   `json:"car"`
	}{}

	err = c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	user := models.User{
		ID:       uint(idNum),
		Phone:    data.Phone,
		RealName: data.Name,
		Password: data.Pwd,
		Role:     data.Role,
		CarID:    data.Car,
	}
	//fmt.Printf("%#v",user)
	err = models.UserUpdateById(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "修改成功",
	})

}

//initUser 初始化用户
func UserInitUser(c *gin.Context) {
	//查看初始化文件是否存在
	initFile := "./assets/init.txt"
	if _, err := os.Stat(initFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("文件不存在")
			//初始化文件目录
			if err := os.Mkdir("./uploads", os.ModePerm); err != nil {
				logs.Error(err)
			}
			if err = os.Mkdir("./logs", os.ModePerm); err != nil {
				logs.Error(err)
			}
			//生成默认车号
			car := models.Car{
				Car: "默认车号",
			}
			if err = models.CarAdd(&car); err != nil {
				logs.Error(err)
			}
			if car, err = models.CarRandomGetOne(); err != nil {
				logs.Error(err)
			}

			//生成超级用户
			phone := utils.RandPhoneStringRunes(11)
			password := utils.RandStringRunes(6)
			user := models.User{
				Phone:    phone,
				RealName: "默认超级管理员",
				Password: password,
				Role:     "admin",
				CarID:    car.ID,
			}
			if _, err := models.UserAddUser(&user); err != nil {
				logs.Error(err)
				c.JSON(http.StatusOK, gin.H{
					"code": 4000,
					"msg":  "系统初始化失败#2",
				})
				return
			}

			logs.Info("初始化账户:%s ||  初始化密码:%s", phone, password)

			//新建init文件
			f, _ := os.Create(initFile)
			defer f.Close()

			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"msg":  "系统已初始化",
				"data": gin.H{
					"phone":    phone,
					"password": password,
				},
			})

			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 4000,
				"msg":  "系统初始化失败#1",
				"data": gin.H{},
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "系统已初始化",
			"data": gin.H{},
		})
		return
	}
}

// UserAddUser 添加用户
func UserAddUser(c *gin.Context) {
	data := struct {
		Phone string `json:"phone"`
		Name  string `json:"name"`
		Role  string `json:"role"`
		Pwd   string `json:"pwd"`
		Car   uint   `json:"car"`
	}{}

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}
	//fmt.Println(data)
	user := models.User{
		Phone:    data.Phone,
		RealName: data.Name,
		Password: data.Pwd,
		Role:     data.Role,
		CarID:    data.Car,
	}

	id, err := models.UserAddUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "添加成功",
		"data": gin.H{
			"id": id,
		},
	})

}

// UserDeleteById 删除用户
func UserDeleteById(c *gin.Context) {
	idStr := c.Param("id")

	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	isDeleted, err := models.UserDelById(uint(idNum))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5000,
			"msg":  err.Error(),
		})
		return
	}

	if isDeleted {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "删除成功",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 4000,
		"msg":  "删除失败",
	})
	return

}

// UserCheckLogin 检查是否登录
func UserCheckLogin(c *gin.Context) {
	jwt := c.GetHeader("jwt")
	logs.Info(jwt)
	phone, _, err := utils.ParseJWT(jwt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4030,
			"msg":  "用户未登录",
		})
		return
	}

	user, err := models.UserGetByPhone(phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "用户已登录",
		"data": gin.H{
			"realname": user.RealName,
			"role":     user.Role,
			"phone":    phone,
			"car":      user.CarID,
			"user":     user,
		},
	})

}

// UserGetAllUser 获取所有用户
func UserGetAllUser(c *gin.Context) {
	data := struct {
		Page    int `json:"page" `
		PerPage int `json:"per_page" `
	}{}

	err := c.BindQuery(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//校验参数
	if data.Page < 1 {
		data.Page = 1
	}

	users, count, err := models.UserGetAllByPage(data.Page, data.PerPage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "ok",
		"data": gin.H{
			"users": users,
			"count": count,
			"page":  data.Page,
		},
	})
}

// UserGetOneUserById 获取单个用户
func UserGetOneUserById(c *gin.Context) {
	idStr := c.Param("id")

	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	user, err := models.UserGetById(uint(idNum))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "ok",
		"data": gin.H{
			"user": user,
		},
	})
}
