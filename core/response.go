package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	Success    int = 200
	Fail       int = 400
	BadRequest int = 401
)

//向客户端发送socket消息结构体
type SocketMessage struct {
	Type   string      `json:"type"`
	Result interface{} `json:"result"`
}

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

//通过websocket向客户端写消息
func ResponseSocketMessage(conn *websocket.Conn, msgType string, result interface{}) {
	var (
		msg = &SocketMessage{
			Type:   msgType,
			Result: result,
		}
		b   []byte
		err error
	)

	if b, err = json.Marshal(msg); err != nil {
		return
	}

	conn.WriteMessage(1, b)
}
