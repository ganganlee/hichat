package route

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/webServer/controller"
	"hichat.zozoo.net/apps/webServer/service"
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
		//根据uuid查找用户
		userRoute.GET("/findByUuid",userCtl.FindByUuid)
	}
}
