package service

import (
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
	return &GatewayService{}
}

//发送消息到网关
func (g *GatewayService) SendMsg(res *SendMsgRequest) (err error) {
	var (
		msg       *models.Message
		tableName string
		redisKey  string
		mqHost    string //mq服务器host
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
	redisKey = "hichat_user:loginStatus:uuid:" + res.ToId + ":string"
	if mqHost = core.CLusterClient.Get(redisKey).Val(); mqHost == "" {
		//未登录，直接返回
		return nil
	}

	//用户是登录状态，可以发送mq消息
	return err
}
