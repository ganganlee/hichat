package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hichat.zozoo.net/apps/messageServer/service"
	"hichat.zozoo.net/core"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ListenController struct {
	service *service.ListenService
}

func NewListenController(s *service.ListenService) *ListenController {
	return &ListenController{
		s,
	}
}

//建立用户长连接
func (l *ListenController) Listen(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		core.ResponseError(c, err.Error())
		return
	}

	l.service.Listen(c.GetString("uuid"), ws)
}
