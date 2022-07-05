package router

import (
	"cailiao_server/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

//设置跨域头
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, jwt")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		//c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// DefineRouter 自定义路由
func DefineRouter(r *gin.Engine) {
	//跨域
	r.Use(cors())

	v1 := r.Group("/v1")
	{
		v1.GET("/", Permission("user"), controllers.Hellow)

		//登录
		v1.POST("/login", controllers.UserLogin)            //用户登录
		v1.GET("/login/status", controllers.UserCheckLogin) //检查登录状态

		//	用户接口
		v1.GET("/users", Permission("admin"), controllers.UserGetAllUser)         //获取所有用户
		v1.POST("/users", Permission("admin"), controllers.UserAddUser)           //添加用户
		v1.DELETE("/users/:id", Permission("admin"), controllers.UserDeleteById)  //删除用户
		v1.GET("/users/:id", Permission("admin"), controllers.UserGetOneUserById) //获取单个用户
		v1.PUT("/users/:id", Permission("admin"), controllers.UserUpdateById)     //更新单个用户

		// 货架接口
		v1.GET("/places", Permission("editor"), controllers.PlaceGetPlaceByPage) //获取货架位置
		v1.GET("/placesall", Permission("editor"), controllers.PlaceGetAll)      //获取所有货架位置
		v1.DELETE("/places/:id", Permission("editor"), controllers.PlaceDelById) //删除某个货架
		v1.GET("/places/:id", Permission("editor"), controllers.PlaceGetOneById) //获取某个货架
		v1.PUT("/places/:id", Permission("editor"), controllers.PlaceUpdateById) //更新某个货架
		v1.POST("/places", Permission("editor"), controllers.PlaceAdd)           //添加某个货架

		// 车号接口
		v1.GET("/car", Permission("admin"), controllers.CarGetAllByPage)   //分页获取车号位置
		v1.GET("/carall", Permission("editor"), controllers.CarGetAll)     //获取所有车号位置
		v1.DELETE("/car/:id", Permission("admin"), controllers.CarDelById) //删除某个车号
		v1.GET("/car/:id", Permission("admin"), controllers.CarGetOneById) //获取某个车号
		v1.PUT("/car/:id", Permission("admin"), controllers.CarUpdateById) //更新某个车号
		v1.POST("/car", Permission("admin"), controllers.CarAdd)           //添加某个车号

		//	材料接口
		//v1.GET("/material", Permission("editor"), controllers.MaterialGetAllByPage)         //获取材料分页
		v1.POST("/material", Permission("editor"), controllers.MaterialAdd)                 //添加材料
		v1.DELETE("/material/id/:id", Permission("editor"), controllers.MaterialDelOneByID) //删除单个材料
		v1.GET("/material/id/:id", Permission("editor"), controllers.MaterialGetOneById)    //获取单个材料
		v1.PUT("/material/id/:id", Permission("editor"), controllers.MaterialUpdateOneById) //更新单个材料
		v1.GET("/material/s", Permission("editor"), controllers.MaterialSearch)             //搜索材料
		v1.POST("/material/all", Permission("editor"), controllers.MaterialAddAll)          //批量添加材料
		v1.GET("/material/download", Permission("editor"), controllers.MaterialDownload)    //下载材料清单

		//	出入库记录接口
		v1.POST("/record", Permission("editor"), controllers.RecordAdd)                         //添加记录
		v1.GET("/record", Permission("editor"), controllers.RecordGetAllByPageAndSearch)        //搜索记录
		v1.GET("/record/id/:id", Permission("editor"), controllers.RecordGetAllByPageAndSearch) //搜索记录
		v1.DELETE("/record/id/:id", Permission("editor"), controllers.RecordDelById)            //删除记录记录

	}

	// 简单的路由组: v2
	// v2 := r.Group("/v2")
	// {
	// 	v2.POST("/login", loginEndpoint)
	// 	v2.POST("/submit", submitEndpoint)
	// 	v2.POST("/read", readEndpoint)
	// }
}
