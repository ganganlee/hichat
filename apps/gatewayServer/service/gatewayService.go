package service

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"hichat.zozoo.net/apps/gatewayServer/common"
	"hichat.zozoo.net/apps/gatewayServer/models"
	"hichat.zozoo.net/core"
)

type (
	GatewayService struct {
		model *models.MessageModel
	}

	//发送消息
	SendMsgRequest struct {
		FromId      string `json:"from_id" validate:"required"`
		ToId        string `json:"to_id" validate:"required"`
		MsgType     string `json:"msg_type" validate:"required"`
		ContentType string `json:"content_type" validate:"required"`
		Content     string `json:"content" validate:"required"`
	}
)

func NewGatewayService(m *models.MessageModel) *GatewayService {
	return &GatewayService{
		m,
	}
}

//发送消息到网关
func (g *GatewayService) SendMsg(res *SendMsgRequest) (err error) {
	var (
		msg       *models.Message
		tableName string

		redisKey = "user:mqHost:uuid:" + res.ToId + ":string:"
		mqHost   = core.CLusterClient.Get(redisKey).Val()
	)

	msg = &models.Message{
		FromId:      res.FromId,
		ToId:        res.ToId,
		MsgType:     res.MsgType,
		ContentType: res.ContentType,
		Content:     res.Content,
	}

	//获取表名称
	tableName = common.GetMessageTable(res.FromId, res.ToId, res.MsgType)

	//将消息入库
	if err = g.model.Create(tableName, msg); err != nil {
		return err
	}

	//判断消息类型，当消息为私聊时往下执行，当消息为群聊时需要获取群成员，逻辑待定
	if res.MsgType == "groupMessage" {
		return nil
	}

	//判断用户是否登录，登陆时将消息发送至mq队列中
	redisKey = "user:mqHost:uuid:" + res.ToId + ":string:"
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
