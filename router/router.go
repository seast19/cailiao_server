package router

import (
	"cailiao_server/controllers"

	"github.com/gin-gonic/gin"
)

// DefineRouter 自定义路由
func DefineRouter(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/", controllers.Hellow)

	}

	// 简单的路由组: v2
	// v2 := r.Group("/v2")
	// {
	// 	v2.POST("/login", loginEndpoint)
	// 	v2.POST("/submit", submitEndpoint)
	// 	v2.POST("/read", readEndpoint)
	// }
}
