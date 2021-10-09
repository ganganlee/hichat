package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"hichat.zozoo.net/apps/messageServer/common"
	"hichat.zozoo.net/core"
	gateway "hichat.zozoo.net/rpc/Gateway"
	"hichat.zozoo.net/rpc/userGroupMembers"
	"time"
)

//消息服务

type (
	MessageService struct {
		conn       *websocket.Conn
		uuid       string
		gatewayRpc gateway.GatewayService
	}

	SendMsgRequest struct {
		Id          string `json:"id" validate:"required"`           //接收id，私聊表示接收用户id，群聊表示群id
		MsgType     string `json:"msg_type" validate:"required"`     //私聊还是群聊
		ContentType string `json:"content_type" validate:"required"` //消息类型，文本还是其他类型
		Content     string `json:"content" validate:"required"`      //消息内容
		FromId      string `json:"from_id"`                          //发送的id
		GroupId     string `json:"group_id"`                         //群id
	}
)

func NewMessageService(conn *websocket.Conn, uuid string) *MessageService {

	var (
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
			micro.WrapClient(
				hystrix.NewClientWrapper(),
				),
		)

		gatewayRpc = gateway.NewGatewayService(common.AppCfg.RpcServer.GatewayRpc, service.Client())
	)

	return &MessageService{
		conn:       conn,
		uuid:       uuid,
		gatewayRpc: gatewayRpc,
	}
}

//客户端发送消息
func (m *MessageService) SendMsg(data string) {
	var (
		err      error
		res      *SendMsgRequest
		validate = validator.New()
	)

	res = new(SendMsgRequest)

	//将消息转换为对象
	if err = json.Unmarshal([]byte(data), res); err != nil {
		core.ResponseSocketMessage(m.conn, "err", err.Error())
		return
	}

	//验证数据
	if err = validate.Struct(res); err != nil {
		core.ResponseSocketMessage(m.conn, "err", err.Error())
		return
	}

	//处理用户缓存
	go m.handleCache(res)

	//向gateway服务器发送消息
	go m.sendMsgToGateway(res)

	core.ResponseSocketMessage(m.conn, "SendStatus", "ok")
}

//增加用户消息缓存
func (m *MessageService) handleCache(msg *SendMsgRequest) {
	var (
		//err            error
		historyService = NewHistoryRecord(m.conn, m.uuid)
		historyMsg     *HistoryMessage
		t              = time.Now()
		redisKey       = "historyRecord:uuid:" + m.uuid + ":messageUser:string"
	)

	//组织消息数据
	historyMsg = &HistoryMessage{
		Id:          msg.Id,                          //好友id
		MessageType: msg.MsgType,                     //消息类型
		ContentType: msg.ContentType,                 //消息内容类型
		Content:     msg.Content,                     //消息内容
		Date:        t.Format("2006-01-02 15:04:05"), //时间
		Uuid:        m.uuid,                          //发送消息的用户uuid
	}

	//缓存当前用户的聊天对象
	core.CLusterClient.Set(redisKey, msg.Id, 5*time.Minute)

	//判断消息类型，进行对应的操作
	switch historyMsg.MessageType {
	case "privateMessage": //私聊
		//添加缓存逻辑，需要修改自己和对方两处缓存，
		//添加自己的不需要添加未读数量，对方需要添加未读数量
		historyService.PushHistoryRecord(m.uuid, msg.Id, historyMsg, false)

		//获取对方的当前的聊天对象
		var (
			messageId string
			unread    = true
		)

		//判断对方当前聊天对象是自己，则不用增加未读消息
		redisKey = "historyRecord:uuid:" + msg.Id + ":messageUser:string"
		if messageId = core.CLusterClient.Get(redisKey).Val(); messageId == m.uuid {
			unread = false
		}

		//添加对方需要添加未读数量
		historyMsg.Id = m.uuid
		historyService.PushHistoryRecord(msg.Id, m.uuid, historyMsg, unread)
		break
	case "groupMessage": //群聊
		//需要获取所有群成员，给对应群成员添加缓存
		var (
			groupMemberService = NewUserGroupMembersService(m.conn, m.uuid)
			list               *userGroupMembers.MembersResponse
			err                error
			messageId          string
		)

		//获取群成员列表
		if list, err = groupMemberService.GetGroupMembers(msg.Id); err != nil {
			core.ResponseSocketMessage(m.conn, "err", err.Error())
			return
		}

		//循环列表，添加缓存
		for _, item := range list.Members {

			//判断是否添加未读消息
			var addUnread = true
			if item.Uuid == m.uuid { //当用户是自己时不需要添加未读消息
				addUnread = false
			}

			//获取当前成员的当前的聊天对象
			redisKey = "historyRecord:uuid:" + item.Uuid + ":messageUser:string"
			if messageId = core.CLusterClient.Get(redisKey).Val(); messageId == msg.Id {
				addUnread = false
			}

			historyService.PushHistoryRecord(item.Uuid, msg.Id, historyMsg, addUnread)
		}
		break
	}
}

//将消息发送至gateway服务器
func (m *MessageService) sendMsgToGateway(res *SendMsgRequest) {

	var (
		rpcRes *gateway.SendMsgRequest
		err    error
	)

	//组织请求rpc方法参数
	rpcRes = &gateway.SendMsgRequest{
		FromId:      m.uuid,
		ToId:        res.Id,
		MsgType:     res.MsgType,
		ContentType: res.ContentType,
		Content:     res.Content,
	}

	if _, err = m.gatewayRpc.SendMsg(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(m.conn, "err", core.DecodeRpcErr(err.Error()).Error())
		return
	}
}
