package router

import (
	"cailiao_server/models"
	"cailiao_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"time"
)

//admin 权限
func Permission(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取jwt
		jwt := c.GetHeader("jwt")

		//解析jwt
		phone, err := utils.ParseJWT(jwt)
		if err != nil {
			//c.AbortWithStatus(http.StatusUnauthorized)

			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 4030,
				"msg":  "用户未登录",
			})
			return
		}

		//检查登录，获取role
		user, err := models.UserGetUserByPhone(phone)
		if err != nil {
			//c.AbortWithStatus(http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 4030,
				"msg":  "用户未登陆",
			})
			return
		}
		//fmt.Printf("用户角色: %s\n", user.Role)

		switch role {
		case "user":
			if user.Role == "user" || user.Role == "editor" || user.Role == "admin" {
				c.Next()
			} else {
				//c.AbortWithStatus(http.StatusForbidden)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}

		case "editor":
			if user.Role == "editor" || user.Role == "admin" {
				c.Next()
			} else {
				//c.AbortWithStatus(http.StatusForbidden)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 4030,
					"msg":  "用户权限不足",
				})
			}
		case "admin":
			if user.Role == "admin" {
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

//logrus中间件
// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	//日志文件
	fileName := path.Join("./logs", "111.txt")

	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	//logger.SetFormatter(&logrus.TextFormatter{})
	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
