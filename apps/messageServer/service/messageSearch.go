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
	"hichat.zozoo.net/apps/messageServer/common"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/messageSearch"
)

//搜索聊天记录
type (
	MessageSearch struct {
		conn             *websocket.Conn
		uuid             string
		messageSearchRpc messageSearch.SearchMessageService
	}

	SearchRequest struct {
		Keywords string `json:"keywords" validate:"required"`
		ToId     string `json:"to_id" validate:"required"`
		FromId   string `json:"from_id"`
		Page     uint32 `json:"page"`
		PageSize uint32 `json:"page_size"`
	}
)

func NewMessageSearch(conn *websocket.Conn, uuid string) *MessageSearch {

	//连接rpc服务
	var (
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		messageSearchRpc = messageSearch.NewSearchMessageService(common.AppCfg.RpcServer.SearchRpc, service.Client())
	)

	fmt.Println(common.AppCfg.RpcServer.SearchRpc)

	return &MessageSearch{
		conn:             conn,
		uuid:             uuid,
		messageSearchRpc: messageSearchRpc,
	}
}

//搜索
func (m *MessageSearch) Search(data string) {
	var (
		validate = validator.New()
		err      error
		res      *SearchRequest
		rpcRes   *messageSearch.SearchRequest
		rpcRsp   *messageSearch.SearchResponse
	)

	res = new(SearchRequest)

	if err = json.Unmarshal([]byte(data), res); err != nil {
		core.ResponseSocketMessage(m.conn, "err", err.Error())
		return
	}

	//参数验证
	if err = validate.Struct(res); err != nil {
		core.ResponseSocketMessage(m.conn, "err", err.Error())
		return
	}

	//调用rpc方法
	rpcRes = &messageSearch.SearchRequest{
		Page:     res.Page,
		PageSize: res.PageSize,
		Keywords: res.Keywords,
		FromId:   m.uuid,
		ToId:     res.ToId,
	}
	if rpcRes.PageSize == 0 {
		rpcRes.PageSize = 10
	}
	if rpcRes.Page == 0 {
		rpcRes.Page = 1
	}

	fmt.Println(rpcRes)
	rpcRsp, err = m.messageSearchRpc.Search(context.TODO(), rpcRes)
	if err != nil {
		fmt.Println("调用rpc方法失败")
		core.ResponseSocketMessage(m.conn, "err", core.DecodeRpcErr(err.Error()).Error())
		return
	}

	fmt.Println(rpcRsp)
}
