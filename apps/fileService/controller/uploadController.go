package controller

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/fileService/service"
	"hichat.zozoo.net/core"
	"mime/multipart"
)

type UploadController struct {
	service *service.UploadService
}

func NewUploadController(s *service.UploadService) *UploadController {
	return &UploadController{
		s,
	}
}

func (i *UploadController) Upload(c *gin.Context) {
	var (
		err  error
		file *multipart.FileHeader
		res  *service.UploadResponse
	)

	// 单文件
	if file, err = c.FormFile("file"); err != nil {
		core.ResponseError(c, err.Error())
		return
	}

	if res, err = i.service.Upload(file, c); err != nil {
		core.ResponseError(c, err.Error())
		return
	}

	core.ResponseSuccess(c, res)
}
