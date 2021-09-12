package rpc

import (
	"context"
	"github.com/go-playground/validator/v10"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/apps/userServer/service"
	"hichat.zozoo.net/rpc/userGroups"
)

type (
	UserGroupsRpc struct {
		service *service.UserGroupsService
	}
)

func NewUserGroupsRpc(s *service.UserGroupsService) *UserGroupsRpc {
	return &UserGroupsRpc{
		s,
	}
}

//创建群
func (u *UserGroupsRpc) CreateGroup(ctx context.Context, res *userGroups.CreateGroupRequest, rsp *userGroups.CreateGroupResponse) error {
	var (
		userGroup = &model.UserGroups{
			Uuid:        res.Group.Uuid,
			Name:        res.Group.Name,
			Avatar:      res.Group.Avatar,
			Description: res.Group.Description,
		}
		err error
	)

	//调用服务方法
	if err = u.service.CreateGroup(userGroup); err != nil {
		return err
	}

	rsp.Gid = userGroup.Gid
	return nil
}

//删除群
func (u *UserGroupsRpc) DelGroup(ctx context.Context, res *userGroups.DelGroupRequest, rsp *userGroups.DelGroupResponse) error {
	var (
		rpcRes = &service.GroupsRequest{
			Uuid: res.Uuid,
			Gid:  res.Gid,
		}
		err      error
		validate = validator.New()
	)

	//参数验证
	if err = validate.Struct(rpcRes); err != nil {
		return err
	}

	//调用服务器方法
	if err = u.service.DelGroup(rpcRes); err != nil {
		return err
	}

	rsp.Msg = "删除成功"
	return nil
}

//获取群列表
func (u *UserGroupsRpc) Groups(ctx context.Context, res *userGroups.GroupsRequest, rsp *userGroups.GroupsResponse) error {
	return nil
}
