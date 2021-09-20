package route

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/fileService/common"
	"hichat.zozoo.net/apps/fileService/controller"
	"hichat.zozoo.net/apps/fileService/models"
	"hichat.zozoo.net/apps/fileService/service"
)

func InitRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	upload := v1.Group("/upload")
	{
		m := models.NewFilesModel(common.AppOrm)
		s := service.NewUploadService(m)
		c := controller.NewUploadController(s)
		upload.POST("", c.Upload)
	}
}
