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
	"hichat.zozoo.net/rpc/userGroupMembers"
)

type (
	UserGroupMembersService struct {
		conn       *websocket.Conn
		uuid       string
		membersRpc userGroupMembers.UserGroupMembersService
	}

	//删除群成员结构体
	RemoveMemberRequest struct {
		Uuid string `json:"uuid" validate:"required"`
		Gid  string `json:"gid" validate:"required"`
	}
)

func NewUserGroupMembersService(conn *websocket.Conn, uuid string) *UserGroupMembersService {
	var (
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		memberRpc = userGroupMembers.NewUserGroupMembersService(common.AppCfg.RpcServer.UserRpc, service.Client())
	)

	return &UserGroupMembersService{
		conn:       conn,
		uuid:       uuid,
		membersRpc: memberRpc,
	}
}

//添加群成员
func (u *UserGroupMembersService) AddMember(str string) {

	var (
		addRequest *RemoveMemberRequest
		res        *userGroupMembers.AddMemberRequest
		err        error
		rsp        *userGroupMembers.AddMemberResponse
		validate   = validator.New()
	)

	addRequest = new(RemoveMemberRequest)
	if err = json.Unmarshal([]byte(str), addRequest); err != nil {
		core.ResponseSocketMessage(u.conn, "err", err.Error())
		return
	}
	if err = validate.Struct(addRequest); err != nil {
		core.ResponseSocketMessage(u.conn, "err", err.Error())
		return
	}

	//赋值请求参数
	res = &userGroupMembers.AddMemberRequest{
		Gid:  addRequest.Gid,
		Uuid: addRequest.Uuid,
	}

	if rsp, err = u.membersRpc.AddMember(context.TODO(), res); err != nil {
		core.ResponseSocketMessage(u.conn, "err", core.DecodeRpcErr(err.Error()).Error())
		return
	}

	core.ResponseSocketMessage(u.conn, "AddMember", rsp.Msg)
}

//删除群成员
func (u *UserGroupMembersService) RemoveMember(data string) {

	var (
		res       *userGroupMembers.DelByMemberIdRequest
		err       error
		rsp       *userGroupMembers.DelByMemberIdResponse
		removeRes *RemoveMemberRequest
		validate = validator.New()
	)

	//将字符串转换为对象
	removeRes = new(RemoveMemberRequest)
	if err = json.Unmarshal([]byte(data), removeRes); err != nil {
		core.ResponseSocketMessage(u.conn, "err", err.Error())
		return
	}

	//验证数据
	if err = validate.Struct(removeRes);err != nil {
		core.ResponseSocketMessage(u.conn, "err", err.Error())
		return
	}

	res = &userGroupMembers.DelByMemberIdRequest{
		Gid:  removeRes.Gid,
		Uuid: removeRes.Uuid,
	}

	//调用rpc方法
	if rsp, err = u.membersRpc.DelByMemberId(context.TODO(), res); err != nil {
		core.ResponseSocketMessage(u.conn, "err", core.DecodeRpcErr(err.Error()).Error())
		return
	}

	core.ResponseSocketMessage(u.conn, "RemoveMember", rsp.Msg)
}

//获取群成员
func (u *UserGroupMembersService) GroupMembers(gid string) {
	if gid == "" {
		core.ResponseSocketMessage(u.conn, "err", "群id不能为空")
		return
	}

	var (
		list *userGroupMembers.MembersResponse
		err  error
	)
	if list, err = u.membersRpc.Members(context.TODO(), &userGroupMembers.MembersRequest{
		Gid: gid,
	}); err != nil {
		core.ResponseSocketMessage(u.conn, "err", err.Error())
		return
	}

	core.ResponseSocketMessage(u.conn, "GroupMembers", list.Members)
}
