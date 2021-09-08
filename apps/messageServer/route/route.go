package route

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/messageServer/controller"
	"hichat.zozoo.net/apps/messageServer/service"
	"hichat.zozoo.net/middleware"
)

func InitListenRoute(c *gin.Engine) {
	v1 := c.Group("/v1")

	//使用授权中间件,
	v1.Use(middleware.Auth())

	//用户长连接
	svc := service.NewListenService()
	ctr := controller.NewListenController(svc)
	v1.GET("/listen", ctr.Listen)
}
