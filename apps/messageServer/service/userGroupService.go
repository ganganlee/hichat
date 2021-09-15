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
	"hichat.zozoo.net/rpc/userGroups"
)

type (
	UserGroupService struct {
		Uuid     string                       //用户uuid
		Conn     *websocket.Conn              //用户长连接
		groupRpc userGroups.UserGroupsService //rpc服务
	}

	//创建群请求
	CreateGroupRequest struct {
		Name        string `json:"name" validate:"required,min=3,max=25"`
		Description string `json:"description" validate:"required"`
		Avatar      string `json:"avatar" validate:"required,url"`
	}

	//修改群信息
	EditGroupRequest struct {
		CreateGroupRequest
		Gid string `json:"gid" validate:"required"`
	}
)

//用户群服务

func NewUserGroupService(uuid string, conn *websocket.Conn) *UserGroupService {

	//注册用户rpc服务
	var (
		service = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		groupRpc = userGroups.NewUserGroupsService(common.AppCfg.RpcServer.UserRpc, service.Client())
	)

	return &UserGroupService{
		Uuid:     uuid,
		Conn:     conn,
		groupRpc: groupRpc,
	}
}

//创建群
func (u *UserGroupService) CreateGroup(res string) {
	var (
		param    = new(CreateGroupRequest)
		validate = validator.New()
		err      error
		rpcRes   *userGroups.CreateGroupRequest
		rpcRsp   *userGroups.CreateGroupResponse
	)

	//解析json
	if err = json.Unmarshal([]byte(res), param); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", err.Error())
		return
	}

	//参数验证
	if err = validate.Struct(param); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", err.Error())
		return
	}

	//调用rpc服务接口
	rpcRes = &userGroups.CreateGroupRequest{
		Group: &userGroups.Group{
			Uuid:        u.Uuid,
			Name:        param.Name,
			Description: param.Description,
			Avatar:      param.Avatar,
		},
	}
	if rpcRsp, err = u.groupRpc.CreateGroup(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	//创建成功
	core.ResponseSocketMessage(u.Conn, "createGroup", rpcRsp.Gid)
}

//删除群
func (u *UserGroupService) DelGroup(gid string) {
	var (
		err    error
		rpcRes *userGroups.DelGroupRequest
		rpcRsp *userGroups.DelGroupResponse
	)

	rpcRes = &userGroups.DelGroupRequest{
		Uuid: u.Uuid,
		Gid:  gid,
	}

	//调用rpc服务
	if rpcRsp, err = u.groupRpc.DelGroup(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(u.Conn, "createGroup", rpcRsp.Msg)
}

//获取群列表
func (u *UserGroupService) Groups(gid string) {
	var (
		res *userGroups.GroupsRequest
		rsp *userGroups.GroupsResponse
		err error
	)

	res = &userGroups.GroupsRequest{
		Uuid: u.Uuid,
	}

	if rsp, err = u.groupRpc.Groups(context.TODO(), res); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(u.Conn, "Groups", rsp)
}

//根据gid查找群
func (u *UserGroupService) FindByGid(gid string) {
	var (
		err error
		rsp *userGroups.FindByGidResponse
	)

	if rsp, err = u.groupRpc.FindByGid(context.TODO(), &userGroups.FindByGidRequest{Gid: gid}); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	core.ResponseSocketMessage(u.Conn, "FindByGid", rsp)
}

//修改群信息
func (u *UserGroupService) EditGroup(res string) {
	var (
		param    = new(EditGroupRequest)
		validate = validator.New()
		err      error
		rpcRes   *userGroups.EditGroupRequest
		rpcRsp   *userGroups.EditGroupResponse
	)

	//解析json
	if err = json.Unmarshal([]byte(res), param); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", err.Error())
		return
	}

	//参数验证
	if err = validate.Struct(param); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", err.Error())
		return
	}

	//调用rpc服务接口
	rpcRes = &userGroups.EditGroupRequest{
		Group: &userGroups.Group{
			Uuid:        u.Uuid,
			Name:        param.Name,
			Description: param.Description,
			Avatar:      param.Avatar,
			Gid:         param.Gid,
		},
	}
	if rpcRsp, err = u.groupRpc.EditGroup(context.TODO(), rpcRes); err != nil {
		core.ResponseSocketMessage(u.Conn, "err", core.DecodeRpcErr(err.Error()))
		return
	}

	//修改成功
	core.ResponseSocketMessage(u.Conn, "updateGroup", rpcRsp.Msg)
}
