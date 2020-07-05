package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Hellow(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// 登录

// 修改用户
// 添加用户
// 删除用户

// 检查是否登录
