package service

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/apps/userServer/common"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/user"
	"hichat.zozoo.net/rpc/userGroups"
	"sync"
)

//群成员服务

var wg *sync.WaitGroup

type (
	UserGroupMembersService struct {
		model    *model.UserGroupMembersModel
		userRpc  user.UserService             //用户服务rpc
		groupRpc userGroups.UserGroupsService //用户群服务rpc
	}

	AddMemberRequest struct {
		Gid  string `json:"gid" validate:"required"`
		Uuid string `json:"uuid" validate:"required"`
	}
)

func NewUserGroupMembersService(m *model.UserGroupMembersModel) *UserGroupMembersService {
	var (
		server = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		userRpc      = user.NewUserService(common.AppCfg.ServerName, server.Client())
		userGroupRpc = userGroups.NewUserGroupsService(common.AppCfg.ServerName, server.Client())
	)
	wg = new(sync.WaitGroup)

	return &UserGroupMembersService{
		m,
		userRpc,
		userGroupRpc,
	}
}

//添加群成员
func (u *UserGroupMembersService) AppendMember(res *AddMemberRequest) (err error) {
	var (
		groupRsp      *userGroups.FindByGidResponse
		userRsp       *user.FindByUuidResponse
		groupErr      error
		userErr       error
		member        *model.UserGroupMembers
		membersKey    string
		userGroupsKey string
	)

	//开启携程查询用户信息是否正确
	wg.Add(2)

	//判断群是否存在
	go func() {
		groupRsp, groupErr = u.groupRpc.FindByGid(context.TODO(), &userGroups.FindByGidRequest{Gid: res.Gid})
		wg.Done()
	}()

	//判断用户是否存在
	go func() {
		userRsp, userErr = u.userRpc.FindByUuid(context.TODO(), &user.FindByUuidRequest{
			Uuid: res.Uuid,
		})
		wg.Done()
	}()
	wg.Wait()

	//判断结果
	if groupErr != nil {
		return groupErr
	}
	if userErr != nil {
		return userErr
	}

	//向数据库添加数据
	member = &model.UserGroupMembers{
		GroupId: res.Gid,
		UserId:  res.Uuid,
	}

	if err = u.model.InsertMember(member); err != nil {
		return err
	}

	//添加群成员缓存
	membersKey = "user_group_members:gid:" + res.Gid + ":member"
	core.CLusterClient.SAdd(membersKey, res.Gid)

	//向用户群缓存添加当前群
	userGroupsKey = "user_groups:uuid:" + res.Uuid + ":member"
	core.CLusterClient.SAdd(userGroupsKey, res.Gid)

	return nil
}

//删除群成员
func (u *UserGroupMembersService) RemoveMember(res *AddMemberRequest) (err error) {
	var (
		member        *model.UserGroupMembers
		membersKey    string
		userGroupsKey string
	)

	member = &model.UserGroupMembers{
		GroupId: res.Gid,
		UserId:  res.Uuid,
	}

	if err = u.model.DelByMemberId(member); err != nil {
		return err
	}

	//从群成员中将当前用户移除
	membersKey = "user_group_members:gid:" + res.Gid + ":member"
	core.CLusterClient.SRem(membersKey, res.Gid)

	//向用户群缓存添加当前群
	userGroupsKey = "user_groups:uuid:" + res.Uuid + ":member"
	core.CLusterClient.SRem(userGroupsKey, res.Gid)

	return err
}

//删除群成员
func (u *UserGroupMembersService) Delete(gid string) (err error) {
	if err = u.model.DelByGroupId(gid); err != nil {
		return err
	}

	var (
		membersKey = "user_group_members:gid:" + gid + ":member"
		members    []string
	)
	//获取当前群的所有成员，从成员的群缓存中删除当前群
	members = core.CLusterClient.SMembers(membersKey).Val()
	for _, uuid := range members {
		userGroupsKey := "user_groups:uuid:" + uuid + ":member"
		core.CLusterClient.SRem(userGroupsKey, gid)
	}

	//删除当前群缓存
	core.CLusterClient.Del(membersKey)
	return nil
}

//获取成员群列表
func (u *UserGroupMembersService) MemberGroups(uuid string) (list []string, err error) {
	var (
		userGroupsKey = "user_groups:uuid:" + uuid + ":member"
		exist         int64
		memberGroups  []model.UserGroupMembers
	)

	//判断缓存存在直接返回
	if exist = core.CLusterClient.Exists(userGroupsKey).Val(); exist != 0 {
		//存在缓存
		list = core.CLusterClient.SMembers(userGroupsKey).Val()
		return list, nil
	}

	//不存在缓存，查询数据库
	if memberGroups, err = u.model.MemberGroups(uuid); err != nil {
		return nil, err
	}

	//循环取到的数据，将数据加入到缓存中
	list = make([]string, 0)
	for _, val := range memberGroups {
		core.CLusterClient.SAdd(userGroupsKey, val.GroupId)
		list = append(list, val.GroupId)
	}

	return list, nil
}
