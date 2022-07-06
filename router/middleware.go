package router

import (
	"cailiao_server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Permission 路由权限拦截器，会话角色权限大于等于设置角色才能访问
func Permission(setRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//解析jwt
		jwt := c.GetHeader("jwt")
		_, userRole, err := utils.ParseJWT(jwt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 4030,
				"msg":  "用户未登录",
			})
			return
		}

		//筛选角色
		switch setRole {
		case "user":
			if userRole == "user" || userRole == "editor" || userRole == "admin" {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}
		case "editor":
			if userRole == "editor" || userRole == "admin" {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}
		case "admin":
			if userRole == "admin" {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}
		default:
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 4030,
				"msg":  "用户权限不足",
			})
		}
	}
}
