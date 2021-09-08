package service

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/apps/messageServer/common"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/user"
)

//用户服务
type (
	UserService struct {
		Uuid    string           //用户uuid
		Conn    *websocket.Conn  //用户长连接
		userRpc user.UserService //用户rpc服务
	}
)

func NewUserService(uuid string, conn *websocket.Conn) *UserService {
	var (

		//注册用户rpc服务
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		userRpc = user.NewUserService(common.AppCfg.RpcServer.UserRpc, service.Client())
	)

	return &UserService{
		Uuid:    uuid,
		Conn:    conn,
		userRpc: userRpc,
	}
}

//根据用户名查找用户
func (u *UserService) FindByName(username string) {
	var (
		rpcRes *user.FindByUsernameRequest
		rpcRsp *user.FindByUsernameResponse
		err    error
	)

	//调用rpc方法
	rpcRes = &user.FindByUsernameRequest{
		Username: username,
	}
	if rpcRsp, err = u.userRpc.FindByUsername(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(u.Conn, "findUser", rpcRsp.Users)
}
