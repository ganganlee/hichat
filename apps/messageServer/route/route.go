package route

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/messageServer/controller"
	"hichat.zozoo.net/apps/messageServer/service"
)

func InitWebRoute(c *gin.Engine) {
	v1 := c.Group("/v1")

	//用户长连接
	svc := service.NewListenService()
	ctr := controller.NewListenController(svc)
	v1.GET("/listen",ctr.Listen)
}

