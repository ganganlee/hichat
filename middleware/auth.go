package middleware

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/core"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			id    string
			err   error
		)

		//获取授权token
		token = c.Request.Header.Get("Authorization")
		if token == "" {
			//当token不存在时，获取地址栏token
			token = c.Query("token")
		}

		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			core.ResponseError(c, "Authorization 不能为空")
			c.Abort()
			return
		}
		if id, err = core.ValidateToken(token); err != nil {
			core.ResponseError(c, err.Error())
			c.Abort()
			return
		}

		//将用户uuid注册在上下文中
		c.Set("uuid", id)

		c.Next()
	}
}
