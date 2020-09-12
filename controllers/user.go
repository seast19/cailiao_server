package controllers

import (
	"cailiao_server/models"
	"cailiao_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Hellow(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// 登录
func UserLogin(c *gin.Context) {
	//获取参数
	data := struct {
		Phone string `json:"phone"`
		Pwd   string `json:"pwd"`
	}{}

	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println("用户输入参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//检验账号密码
	isCurrentPwd, err := models.UserCheckPwd(data.Phone, data.Pwd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5000,
			"msg":  err.Error(),
		})
		return
	}

	//账户密码错误
	if !isCurrentPwd {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "账户或密码错误",
		})
		return
	}

	//生成jwt
	jwt, err := utils.GenJWT(data.Phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5000,
			"msg":  "生成jwt失败",
		})
		return
	}

	//获取用户角色
	user, err := models.UserGetUserByPhone(data.Phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5000,
			"msg":  err.Error(),
		})
		return
	}

	//返回数据
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "登录成功",
		"data": gin.H{
			"jwt":   jwt,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}

// 修改用户
func UserUpdateUserById(c *gin.Context) {
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
	}

	err = models.UserUpdateUserById(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "ok",
	})

}

// 添加用户
func UserAddUser(c *gin.Context) {
	data := struct {
		Phone string `json:"phone"`
		Name  string `json:"name"`
		Role  string `json:"role"`
		Pwd   string `json:"pwd"`
	}{}

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	user := models.User{
		Phone:    data.Phone,
		RealName: data.Name,
		Password: data.Pwd,
		Role:     data.Role,
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

// 删除用户
func UserDeleteUserById(c *gin.Context) {
	idStr := c.Param("id")

	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	isDeleted, err := models.UserDelUserById(uint(idNum))
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

// 检查是否登录
func UserCheckLogin(c *gin.Context) {
	jwt := c.GetHeader("jwt")

	phone, err := utils.ParseJWT(jwt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4030,
			"msg":  "用户未登录",
		})
		return
	}

	user, err := models.UserGetUserByPhone(phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "用户已登录",
		"data": gin.H{
			"role":  user.Role,
			"phone": phone,
		},
	})

}

//获取所有用户
func UserGetAllUser(c *gin.Context) {
	data := struct {
		Page    int `json:"page" form:"page"`
		PerPage int `json:"per_page" form:"per_page"`
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
	//if data.PerPage < 10 || data.PerPage > 100 {
	//	data.PerPage = 10
	//}

	users, count, err := models.UserGetUsersByPage(data.Page, data.PerPage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
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

//获取单个用户
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

	user, err := models.UserGetUserById(uint(idNum))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
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
