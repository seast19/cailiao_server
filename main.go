package main

import (
	"cailiao_server/config"
	"cailiao_server/router"
	"cailiao_server/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//设置模式
	gin.SetMode(config.RUN_MOD)

	//设置日志
	utils.SetLogger()

	// 注册路由
	router.DefineRouter(r)
	r.Static("/statics", "./statics")

	logs.Info("系统启动成功")

	_ = r.Run(":9503")
}
