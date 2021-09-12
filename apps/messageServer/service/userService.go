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
	"hichat.zozoo.net/rpc/userFriends"
)

//用户服务
type (
	UserService struct {
		Uuid           string                         //用户uuid
		Conn           *websocket.Conn                //用户长连接
		userRpc        user.UserService               //用户rpc服务
		userFriendsRpc userFriends.UserFriendsService //用户好友rpc服务
	}
)

func NewUserService(uuid string, conn *websocket.Conn) *UserService {
	var (
		//注册用户rpc服务
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		userRpc = user.NewUserService(common.AppCfg.RpcServer.UserRpc, service.Client())

		//注册用户好友服务
		userFriendsRpc = userFriends.NewUserFriendsService(common.AppCfg.RpcServer.UserRpc, service.Client())
	)

	return &UserService{
		Uuid:           uuid,
		Conn:           conn,
		userRpc:        userRpc,
		userFriendsRpc: userFriendsRpc,
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

//添加好友申请
func (u *UserService) ApplyFriend(friendUuid string) {
	var (
		err    error
		rpcRes = &userFriends.ApplyFriendsRequest{
			Uuid:       u.Uuid,
			FriendUuid: friendUuid,
		}
		rpcRsp *userFriends.ApplyFriendsResponse
	)

	//调用rpc方法
	if rpcRsp, err = u.userFriendsRpc.ApplyFriends(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", err.Error())
		return
	}

	//添加好友发送成功
	core.ResponseSocketMessage(u.Conn, "success", rpcRsp.Msg)
}

//同意添加好友
func (u *UserService) ApproveFriend(friendUuid string) {
	var (
		err    error
		rpcRes = &userFriends.ApproveFriendsRequest{
			Uuid:       u.Uuid,
			FriendUuid: friendUuid,
		}
		rpcRsp *userFriends.ApproveFriendsResponse
	)

	//调用rpc方法
	if rpcRsp, err = u.userFriendsRpc.ApproveFriends(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(u.Conn, "success", rpcRsp.Msg)

}

//拒绝好友申请
func (u *UserService) RefuseFriend(friendUuid string) {
	var (
		err    error
		rpcRes = &userFriends.RefuseFriendsRequest{
			Uuid:       u.Uuid,
			FriendUuid: friendUuid,
		}
		rpcRsp *userFriends.RefuseFriendsResponse
	)

	//调用rpc方法
	if rpcRsp, err = u.userFriendsRpc.RefuseFriends(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(u.Conn, "success", rpcRsp.Msg)
}

//删除好友
func (u *UserService) DelFriend(friendUuid string) {
	var (
		err    error
		rpcRes = &userFriends.DelFriendsRequest{
			Uuid:       u.Uuid,
			FriendUuid: friendUuid,
		}

		rpcRsp *userFriends.DelFriendsResponse
	)

	//调用rpc方法
	if rpcRsp, err = u.userFriendsRpc.DelFriends(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(u.Conn, "success", rpcRsp.Msg)
}

//获取好友列表
func (u *UserService) Friends(val string) {
	var (
		err    error
		rpcRes = &userFriends.FriendsRequest{
			Uuid: u.Uuid,
		}
		rpcRsp *userFriends.FriendsResponse
	)

	//调用rpc方法
	if rpcRsp, err = u.userFriendsRpc.Friends(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	//返回好友列表
	core.ResponseSocketMessage(u.Conn, "friends", rpcRsp.Friends)
}
