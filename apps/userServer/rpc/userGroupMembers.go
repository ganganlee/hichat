package rpc

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"hichat.zozoo.net/apps/userServer/service"
	"hichat.zozoo.net/rpc/userGroupMembers"
)

//群成员管理rpc服务

type (
	GroupMembersRpc struct {
		service *service.UserGroupMembersService
	}
)

func NewGroupMembersRpc(s *service.UserGroupMembersService) *GroupMembersRpc {
	return &GroupMembersRpc{
		s,
	}
}

//添加群成员
func (u *GroupMembersRpc) AddMember(ctx context.Context, res *userGroupMembers.AddMemberRequest, rsp *userGroupMembers.AddMemberResponse) error {
	var (
		err      error
		rpcRes   *service.AddMemberRequest
		validate = validator.New()
	)

	rpcRes = &service.AddMemberRequest{
		Gid:  res.Gid,
		Uuid: res.Uuid,
	}
	if err = validate.Struct(rpcRes); err != nil {
		return err
	}

	if err = u.service.AppendMember(rpcRes); err != nil {
		return err
	}

	rsp.Msg = "添加成功"
	return nil
}

//删除群成员
func (u *GroupMembersRpc) DelByMemberId(ctx context.Context, res *userGroupMembers.DelByMemberIdRequest, rsp *userGroupMembers.DelByMemberIdResponse) error {
	var (
		param    *service.AddMemberRequest
		err      error
		validate = validator.New()
	)

	param = &service.AddMemberRequest{
		Gid:  res.Gid,
		Uuid: res.Uuid,
	}

	//参数验证
	if err = validate.Struct(param); err != nil {
		return err
	}

	//调用服务方法
	if err = u.service.RemoveMember(param); err != nil {
		return err
	}
	rsp.Msg = "删除成功！"
	return nil
}

//解散群
func (u *GroupMembersRpc) DelMembers(ctx context.Context, res *userGroupMembers.DelMembersRequest, rsp *userGroupMembers.DelMembersResponse) error {
	var err error
	if res.Gid == "" {
		return errors.New("群id不能为空")
	}

	if err = u.service.Delete(res.Gid); err != nil {
		return err
	}

	rsp.Msg = "删除成功！"
	return nil
}

//获取成员群列表
func (u *GroupMembersRpc) MemberGroups(ctx context.Context, res *userGroupMembers.MemberGroupsRequest, rsp *userGroupMembers.MemberGroupsResponse) error {
	var (
		err  error
		list []string
	)

	if res.Uuid == "" {
		return errors.New("用户uuid不能为空")
	}

	if list, err = u.service.MemberGroups(res.Uuid); err != nil {
		return err
	}

	rsp.Groups = list
	return nil
}
