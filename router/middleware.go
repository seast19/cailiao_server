package router

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

//
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置 example 变量
		c.Set("example", "12345")

		// 请求前

		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Print(latency)

		// 获取发送的 status
		status := c.Writer.Status()
		log.Println(status)
	}
}
