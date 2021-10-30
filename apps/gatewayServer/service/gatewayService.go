package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/streadway/amqp"
	"hichat.zozoo.net/apps/gatewayServer/common"
	"hichat.zozoo.net/apps/gatewayServer/models"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/userGroupMembers"
)

type (
	GatewayService struct {
		model     *models.MessageModel
		memberRpc userGroupMembers.UserGroupMembersService
	}

	//发送消息
	SendMsgRequest struct {
		FromId      string `json:"from_id" validate:"required"`
		ToId        string `json:"to_id" validate:"required"`
		MsgType     string `json:"msg_type" validate:"required"`
		ContentType string `json:"content_type" validate:"required"`
		Content     string `json:"content" validate:"required"`
		GroupId     string `json:"group_id"`
	}
)

func NewGatewayService(m *models.MessageModel) *GatewayService {
	var (
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		memberRpc = userGroupMembers.NewUserGroupMembersService(common.AppCfg.RpcServer.UserRpc, service.Client())
	)

	return &GatewayService{
		model:     m,
		memberRpc: memberRpc,
	}
}

//发送消息到网关
func (g *GatewayService) SendMsg(res *SendMsgRequest) (err error) {

	//判断消息类型，当消息为私聊时往下执行，当消息为群聊时需要获取群成员，逻辑待定
	switch res.MsgType {
	case "groupMessage": //群聊
		err = g.sendGroupMessage(res)
		break
	case "privateMessage": //私聊
		//消息入库
		if err = g.saveMsg(res); err != nil {
			return err
		}
		err = g.sendMq(res.ToId, res)
		break
	case "ApplyFriend": //好友申请通知
		err = g.sendMq(res.ToId, res)
		break
	case "ApproveFriend":
		err = g.sendMq(res.ToId, res)
		break
	case "AddMember"://邀请入群通知
		err = g.sendMq(res.ToId, res)
		break
	case "Refresh"://刷新列表通知
		err = g.sendMq(res.ToId, res)
	}

	return err
}

//发送私聊消息
func (g *GatewayService) sendGroupMessage(msg *SendMsgRequest) (err error) {

	//消息入库
	if err = g.saveMsg(msg); err != nil {
		return err
	}

	var (
		rsp *userGroupMembers.MembersResponse
	)

	//获取群成员
	if rsp, err = g.memberRpc.Members(context.TODO(), &userGroupMembers.MembersRequest{Gid: msg.ToId}); err != nil {
		return err
	}

	msg.GroupId = msg.ToId

	//循环群成员，判断对方是否登录，如果登录则发送websocket消息
	for _, item := range rsp.Members {

		//判断，当当前成员为发送消息的成员时，不需要发送mq消息
		if msg.FromId == item.Uuid {
			continue
		}

		msg.ToId = item.Uuid
		if err = g.sendMq(item.Uuid, msg); err != nil {
			continue
		}
	}
	return nil
}

//发送消息到rabbitMq
func (g *GatewayService) sendMq(id string, msg *SendMsgRequest) (err error) {

	var (
		redisKey string
		mqHost   string
	)

	//判断用户是否登录，登陆时将消息发送至mq队列中
	redisKey = "user:mqHost:uuid:" + id + ":string:"
	if mqHost = core.CLusterClient.Get(redisKey).Val(); mqHost == "" {
		//未登录，直接返回
		return nil
	}

	//将消息发送至rabbitMq
	var (
		mqList = common.MqQueue
		mq     *amqp.Channel
		exist  bool
		b      []byte
	)

	//判断当前mq是否存在
	if mq, exist = mqList[mqHost]; !exist {
		return errors.New(mqHost + "消息队列主机未链接")
	}

	//将对象转为字符切片
	if b, err = json.Marshal(msg); err != nil {
		return err
	}

	//将消息发送至mq队列中
	if err = common.Publish(mq, b); err != nil {
		return err
	}

	return err
}

//将消息保存进入数据库
func (g *GatewayService) saveMsg(res *SendMsgRequest) (err error) {
	var (
		msg       *models.Message
		tableName string
	)

	msg = &models.Message{
		FromId:      res.FromId,
		ToId:        res.ToId,
		MsgType:     res.MsgType,
		ContentType: res.ContentType,
		Content:     res.Content,
	}

	//获取表名称
	tableName = core.GetMessageTable(res.FromId,res.ToId,res.MsgType)

	//将消息入库
	if err = g.model.Create(tableName, msg); err != nil {
		return err
	}

	return nil
}
