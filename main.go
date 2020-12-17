package main

import (
	"cailiao_server/router"
	"cailiao_server/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 记录到文件。
	//f, _ := os.Create("./logs/gin.log")
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//r.Use(router.LoggerToFile())\

	//gin.SetMode(gin.ReleaseMode)

	//设置日志
	utils.SetLogger()

	// 路由
	router.DefineRouter(r)
	// Listen and Server in 0.0.0.0:8080

	utils.Mlogger.Debug("+----- 系统启动成功 -----+")

	r.Run(":7090")
}
