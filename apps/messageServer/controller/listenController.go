package controller

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/messageServer/service"
)

type ListenController struct {
	service *service.ListenService
}

func NewListenController(s *service.ListenService) *ListenController {
	return &ListenController{
		s,
	}
}

func (l *ListenController) Listen(c *gin.Context) {

}
