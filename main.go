package main

import (
	"cailiao_server/router"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 记录到文件。
	f, _ := os.Create("./logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 路由
	router.DefineRouter(r)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
