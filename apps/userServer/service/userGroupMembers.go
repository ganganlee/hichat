package service

import (
	"context"
	"encoding/json"
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

	//群成员信息
	MemberInfo struct {
		Uuid     string
		Username string
		Avatar   string
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
		memberInfo    *MemberInfo
		b             []byte
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

	//需要获取成员信息，调用rpc接口，因为获取群成员时需要使用

	//添加群成员缓存
	membersKey = "user_group_members:gid:" + res.Gid + ":member"

	//组织群成员信息，并且将信息转换为接送
	memberInfo = &MemberInfo{
		Uuid:     userRsp.User.Uuid,
		Username: userRsp.User.Username,
		Avatar:   userRsp.User.Avatar,
	}

	if b, err = json.Marshal(memberInfo); err != nil {
		return err
	}
	core.CLusterClient.HSet(membersKey, memberInfo.Uuid, string(b))

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

	//从群成员中将当前群中移除
	membersKey = "user_group_members:gid:" + res.Gid + ":member"
	core.CLusterClient.HDel(membersKey, res.Uuid)

	//向用户群缓存添移除前群
	userGroupsKey = "user_groups:uuid:" + res.Uuid + ":member"
	core.CLusterClient.SRem(userGroupsKey, res.Gid)

	return err
}

//删除所有群成员
func (u *UserGroupMembersService) Delete(gid string) (err error) {
	if err = u.model.DelByGroupId(gid); err != nil {
		return err
	}

	var (
		membersKey = "user_group_members:gid:" + gid + ":member"
		members    map[string]string
	)
	//获取当前群的所有成员，从成员的群缓存中删除当前群
	members = core.CLusterClient.HGetAll(membersKey).Val()
	for uuid, _ := range members {
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

//获取群所有成员
func (u *UserGroupMembersService) Members(gid string) (m map[string]string, err error) {
	var (
		membersKey = "user_group_members:gid:" + gid + ":member"
		list       []model.UserGroupMembers
		data       map[string]string
	)

	if m = core.CLusterClient.HGetAll(membersKey).Val(); len(m) > 0 {
		return m, err
	}

	//缓存不存在，调用数据库查询
	if list, err = u.model.Members(gid); err != nil {
		return nil, err
	}

	//创建变量接收数据
	data = make(map[string]string, 0)

	//循环列表，将数据加入到缓存中
	for _, val := range list {
		wg.Add(1)

		//并发执行程序
		go func(res model.UserGroupMembers) {
			defer wg.Done()
			//调用rpc方法获取用户信息

			userRsp, userErr := u.userRpc.FindByUuid(context.TODO(), &user.FindByUuidRequest{
				Uuid: res.UserId,
			})
			if userErr != nil {
				return
			}

			//将用户信息加入到群缓存中
			var (
				memberInfo = &MemberInfo{
					Uuid:     userRsp.User.Uuid,
					Username: userRsp.User.Username,
					Avatar:   userRsp.User.Avatar,
				}
				b []byte
			)
			if b, err = json.Marshal(memberInfo); err != nil {
				return
			}

			core.CLusterClient.HSet(membersKey, memberInfo.Uuid, string(b))

			//将信息加入到全局变量中，需要返回
			data[memberInfo.Uuid] = string(b)
		}(val)
	}

	wg.Wait()

	return data, nil
}
