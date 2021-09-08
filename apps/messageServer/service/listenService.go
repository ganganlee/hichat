package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type (
	ListenService struct {
	}

	//通用消息结构体
	ClientMessage struct {
		Type    string `json:"type" validate:"required"`    //消息类型
		Content string `json:"content" validate:"required"` //消息类容
	}
)

//定义全局变量，保存用户连接
var Conns map[string]*websocket.Conn

func NewListenService() *ListenService {
	Conns = make(map[string]*websocket.Conn, 0)

	return &ListenService{}
}

//监听用户长连接
func (l *ListenService) Listen(uuid string, conn *websocket.Conn) (err error) {

	//保存用户连接
	Conns[uuid] = conn

	//开辟协程监听连接状态
	go l.listenStatus(uuid, conn)

	return err
}

//监听客户端连接状态
func (l *ListenService) listenStatus(uuid string, conn *websocket.Conn) {
	var (
		p   []byte
		err error
	)

	defer conn.Close()
	defer delete(Conns, uuid)
	p = make([]byte, 1024)

	for {
		_, p, err = conn.ReadMessage()
		if err != nil {
			break
		}

		//异步处理消息
		l.handleClientMessage(uuid, p)
	}
}

//处理websocket消息
func (l *ListenService) handleClientMessage(uuid string, msg []byte) {
	var (
		clientMessage *ClientMessage
		err           error
	)

	clientMessage = new(ClientMessage)

	err = json.Unmarshal(msg, clientMessage)
	if err != nil {
		return
	}

	//验证消息格式是否正确
	validate := validator.New()
	if err = validate.Struct(clientMessage); err != nil {
		return
	}

	switch clientMessage.Type {
	case "findUser": //根据用户名查找用户
		var userService = NewUserService(uuid, Conns[uuid])
		userService.FindByName(clientMessage.Content)

	case "privateMsg":
		//私聊
		fmt.Println(clientMessage.Content)
		//m.sendPrivateMsg(uuid, clientMessage)
		break
	case "groupMsg":
		//群聊
		//var (
		//	rpcRes GroupMessage.MessageRequest
		//	rsp    *Users.FindByTokenResponse
		//)
		//
		////获取用户信息
		//if rsp, err = m.userRpc.FindByToken(context.TODO(), &Users.FindByTokenRequest{Token: uuid}); err != nil {
		//	return
		//}
		//
		////将消息发送至群聊微服务，通过群聊微服务网关转发值所有用户
		//rpcRes = GroupMessage.MessageRequest{
		//	Token: msgRequest.ToToken,
		//	Body: &GroupMessage.GroupBody{
		//		Type:     msgRequest.Body.Type,
		//		Content:  msgRequest.Body.Content,
		//		Nickname: rsp.User.Username,
		//		HeadImg:  rsp.User.HeadImg,
		//		Token:    uuid,
		//	},
		//}
		//
		//if _, err = m.groupPrc.Send(context.TODO(), &rpcRes); err != nil {
		//	log.Println(err)
		//}
		break
	default:
		fmt.Println(clientMessage.Content)
	}
}
