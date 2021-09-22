package service

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/apps/messageServer/common"
	"hichat.zozoo.net/core"
	gateway "hichat.zozoo.net/rpc/Gateway"
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
	}
)

func NewMessageService(conn *websocket.Conn, uuid string) *MessageService {

	var (
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		gatewayRpc = gateway.NewGatewayService(common.AppCfg.RpcServer.UserRpc, service.Client())
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
		rpcRes   *gateway.SendMsgRequest
		rpcRsp   *gateway.SendMsgResponse
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

	//组织请求rpc方法参数
	rpcRes = &gateway.SendMsgRequest{
		FromId:      res.FromId,
		ToId:        res.Id,
		MsgType:     res.MsgType,
		ContentType: res.ContentType,
		Content:     res.Content,
	}

	if rpcRsp, err = m.gatewayRpc.SendMsg(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(m.conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(m.conn, "SendStatus", rpcRsp.Msg)
}
