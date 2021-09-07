package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Success    int = 200
	Fail       int = 400
	BadRequest int = 401
)

//请求成功
func ResponseSuccess(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":   Success,
		"result": result,
		"msg":    "ok",
	})
}

//请求失败
func ResponseError(c *gin.Context, msg string, code ...int) {

	//判断是否传入code,未传入默认code为fail
	var rspCode = Fail
	if len(code) != 0 {
		rspCode = code[0]
	}

	c.JSON(http.StatusOK, gin.H{
		"code": rspCode,
		"msg":  msg,
	})
}
