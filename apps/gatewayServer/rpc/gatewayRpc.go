package rpc

import (
	"context"
	"github.com/go-playground/validator/v10"
	"hichat.zozoo.net/apps/gatewayServer/service"
	gateway "hichat.zozoo.net/rpc/Gateway"
)

type (
	GatewayRpc struct {
		service *service.GatewayService
	}
)

func NewGatewayRpc(s *service.GatewayService) *GatewayRpc {
	return &GatewayRpc{
		s,
	}
}

//发送消息到网关
func (g *GatewayRpc) SendMsg(ctx context.Context, res *gateway.SendMsgRequest, rsp *gateway.SendMsgResponse) error {
	var (
		err      error
		msg      *service.SendMsgRequest
		validate = validator.New()
	)

	msg = &service.SendMsgRequest{
		FromId:      res.FromId,
		ToId:        res.ToId,
		MsgType:     res.MsgType,
		ContentType: res.ContentType,
		Content:     res.Content,
	}

	//参数验证
	if err = validate.Struct(msg); err != nil {
		return err
	}

	//调用服务方法
	if err = g.service.SendMsg(msg); err != nil {
		return err
	}
	rsp.Msg = "ok"
	return nil
}
