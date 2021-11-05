package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"hichat.zozoo.net/apps/messageServer/common"
	"hichat.zozoo.net/core"
	"log"
	"reflect"
)

type (
	ListenService struct {
	}

	//通用消息结构体
	ClientMessage struct {
		Type     string `json:"type" validate:"required"`     //消息类型
		Services string `json:"services" validate:"required"` //使用的服务
		Content  string `json:"content"`                      //消息类容
	}

	//mq消息结构体
	MqRequest struct {
		FromId      string `json:"from_id"`
		ToId        string `json:"to_id"`
		MsgType     string `json:"msg_type"`
		ContentType string `json:"content_type"`
		Content     string `json:"content"`
		GroupId     string `json:"group_id"`
	}
)

//定义全局变量，保存用户连接
var Conns map[string]*websocket.Conn

func NewListenService() *ListenService {
	Conns = make(map[string]*websocket.Conn, 0)

	//接收mq消息
	var l = &ListenService{}
	go l.receiveMqMsg()

	return l
}

//监听用户长连接
func (l *ListenService) Listen(uuid string, conn *websocket.Conn) (err error) {

	//保存用户连接
	Conns[uuid] = conn

	//开辟协程监听连接状态
	go l.listenStatus(uuid, conn)

	/**
	将本地mq地址保存在用户缓存中，
	标记用户当前是登录状态，
	当其他用户给他发送消息时需要判断当前用户状态，
	当用户是登录状态时，向用户发送socket消息
	*/
	var redisKey = "user:mqHost:uuid:" + uuid + ":string:"
	core.CLusterClient.Set(redisKey, common.AppCfg.MqHost, core.DefaultExpire)

	return err
}

//监听客户端连接状态
func (l *ListenService) listenStatus(uuid string, conn *websocket.Conn) {
	var (
		p   []byte
		err error
	)

	//当用户断开连接时执行
	defer func() {
		//关闭用户socket连接
		conn.Close()
		//删除用户连接信息
		delete(Conns, uuid)
		//删除用户登录状态
		var redisKey = "user:mqHost:uuid:" + uuid + ":string:"
		core.CLusterClient.Del(redisKey)
	}()

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
		core.ResponseSocketMessage(Conns[uuid], "err", err.Error())
		return
	}

	//验证消息格式是否正确
	validate := validator.New()
	if err = validate.Struct(clientMessage); err != nil {
		core.ResponseSocketMessage(Conns[uuid], "err", err.Error())
		return
	}

	//获取服务对象
	var (
		s interface{}
		f []reflect.Value
	)
	if s, err = l.getServiceInterface(clientMessage.Services, uuid); err != nil {
		core.ResponseSocketMessage(Conns[uuid], "err", err.Error())
		return
	}

	//调用服务方法
	if f, err = core.CallFuncByName(s, clientMessage.Type, clientMessage.Content); err != nil {
		core.ResponseSocketMessage(Conns[uuid], "err", "方法"+clientMessage.Type+"不存在")
		return
	}
	_ = f
}

//获取服务对象
func (l *ListenService) getServiceInterface(services string, uuid string) (s interface{}, err error) {
	switch services {
	case "UserService": //用户相关服务
		s = NewUserService(uuid, Conns[uuid])
		break
	case "UserGroupsService": //用户群相关
		s = NewUserGroupService(uuid, Conns[uuid])
		break
	case "UserGroupMemberService": //用户群成员相关
		s = NewUserGroupMembersService(Conns[uuid], uuid)
		break
	case "HistoryRecordService": //用户历史消息相关
		s = NewHistoryRecord(Conns[uuid], uuid)
		break
	case "messageService":
		s = NewMessageService(Conns[uuid], uuid)
		break
	case "messageSearchService": //用户搜索消息
		s = NewMessageSearch(Conns[uuid], uuid)
		break
	default:
		return nil, errors.New("服务不存在")
	}
	return s, nil
}

//接收rabbitMq消息
func (l *ListenService) receiveMqMsg() {

	var ch = common.MqCh
	msgs, err := ch.Consume(
		common.QueueName,
		"MsgWorkConsumer",
		false, //Auto Ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		//将数据发送至具体方法处理
		go l.handleMqMsg(msg.Body)

		//Ack
		if err = msg.Ack(false); err != nil {
			fmt.Println(err)
		}
	}
}

func (l *ListenService) handleMqMsg(msg []byte) {
	var (
		res   *MqRequest
		err   error
		conn  *websocket.Conn
		exist bool
	)

	//将字符切片解析为结构体对象
	res = new(MqRequest)
	if err = json.Unmarshal(msg, res); err != nil {
		fmt.Println(string(msg))
		return
	}

	//判断用户是否登录
	if conn, exist = Conns[res.ToId]; !exist {
		//用户未登录，删除用户缓存
		var redisKey = "user:mqHost:uuid:" + res.ToId + ":string:"
		core.CLusterClient.Del(redisKey)
		return
	}

	//向用户发送消息
	core.ResponseSocketMessage(conn, "MqMsg", res)
}
