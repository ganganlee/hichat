package route

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/webServer/controller"
	"hichat.zozoo.net/apps/webServer/service"
	"hichat.zozoo.net/middleware"
)

func InitWebRoute(c *gin.Engine) {
	v1 := c.Group("/v1")

	//用户相关路由
	userRoute := v1.Group("/user")
	{
		userSvc := service.NreUserService()
		userCtl := controller.NewUserController(userSvc)

		//用户注册
		userRoute.POST("/register", userCtl.Register)
		//用户登录
		userRoute.POST("/login", userCtl.Login)

		//获取用户信息
		userRoute.GET("/info", middleware.Auth(), userCtl.FindByUuid)
	}
}
