package router

import (
	"cailiao_server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

//admin 权限
func Permission(targetRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取jwt
		jwt := c.GetHeader("jwt")

		//解析jwt
		_, role, err := utils.ParseJWT(jwt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 4030,
				"msg":  "用户未登录",
			})
			return
		}

		//检查登录，获取role
		//user, err := models.UserGetUserByPhone(phone)
		//if err != nil {
		//	//c.AbortWithStatus(http.StatusUnauthorized)
		//	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		//		"code": 4030,
		//		"msg":  "用户未登陆",
		//	})
		//	return
		//}
		//fmt.Printf("用户角色: %s\n", user.Role)

		switch targetRole {
		case "user":
			if role == "user" || role == "editor" || role == "admin" {
				c.Next()
			} else {
				//c.AbortWithStatus(http.StatusForbidden)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}

		case "editor":
			if role == "editor" || role == "admin" {
				c.Next()
			} else {
				//c.AbortWithStatus(http.StatusForbidden)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}
		case "admin":
			if role == "admin" {
				c.Next()
			} else {
				//c.AbortWithStatus(http.StatusForbidden)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}
		default:
			//c.AbortWithStatus(http.StatusForbidden)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 4030,
				"msg":  "用户权限不足",
			})
		}
	}
}
