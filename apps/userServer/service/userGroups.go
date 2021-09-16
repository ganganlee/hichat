package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/apps/userServer/common"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/userGroupMembers"
)

type (
	UserGroupsService struct {
		model           *model.UserGroupsModel
		groupMembersRpc userGroupMembers.UserGroupMembersService
	}

	//获取群列表结构体
	GroupsRequest struct {
		Gid  string `json:"gid" validate:"required"`
		Uuid string `json:"uuid" validate:"required"`
	}

	//修改群信息请求
	EditGroupRequest struct {
		Name        string `json:"name" validate:"required,min=3,max=6"`
		Description string `json:"description" validate:"required"`
		Avatar      string `json:"avatar" validate:"required,url"`
		Uuid        string `json:"uuid" validate:"required"`
		Gid         string `json:"gid" validate:"required"`
	}
)

func NewUserGroupsService(m *model.UserGroupsModel) *UserGroupsService {
	var (
		server = micro.NewService(
			micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
		)
		groupMembersRpc = userGroupMembers.NewUserGroupMembersService(common.AppCfg.ServerName, server.Client())
	)
	return &UserGroupsService{
		m,
		groupMembersRpc,
	}
}

//创建群
func (u *UserGroupsService) CreateGroup(g *model.UserGroups) (err error) {
	//创建用户gid
	g.Gid = uuid.New().String()

	//调用群成员rpc方法添加成员
	if err = u.model.Create(g); err != nil {
		return err
	}

	//调用群成员rpc方法添加成员
	var res = &userGroupMembers.AddMemberRequest{
		Gid:  g.Gid,
		Uuid: g.Uuid,
	}
	if _, err = u.groupMembersRpc.AddMember(context.TODO(), res); err != nil {
		return err
	}

	return nil
}

//删除群
func (u *UserGroupsService) DelGroup(res *GroupsRequest) (err error) {
	var (
		rpcRes = &model.UserGroups{
			Uuid: res.Uuid,
			Gid:  res.Gid,
		}
	)

	//TODO 删除缓存
	return u.model.DelByModel(rpcRes)
}

//获取群列表
func (u *UserGroupsService) Groups(uuid string) (list []model.UserGroups, err error) {
	var (
		rpcRsp *userGroupMembers.MemberGroupsResponse
	)
	//调用rpc方法获取群id列表
	if rpcRsp, err = u.groupMembersRpc.MemberGroups(context.TODO(), &userGroupMembers.MemberGroupsRequest{Uuid: uuid}); err != nil {
		return nil, err
	}

	list = make([]model.UserGroups, 0)
	for _, val := range rpcRsp.Groups {
		var (
			m *model.UserGroups
		)
		if m, err = u.FindByGid(val); err != nil {
			return nil, err
		}

		list = append(list, model.UserGroups{
			Gid:         m.Gid,
			Uuid:        m.Uuid,
			Avatar:      m.Avatar,
			Description: m.Description,
			Name:        m.Name,
		})
	}

	return list, nil
}

//获取群信息，根据群id
func (u *UserGroupsService) FindByGid(gid string) (userGroup *model.UserGroups, err error) {
	var (
		redisKey = "user_groups:gid:" + gid + ":json" //缓存key
		b        []byte                               //字符切片
	)

	userGroup = new(model.UserGroups)

	//获取缓存
	if b, err = core.CLusterClient.Get(redisKey).Bytes(); err == nil {
		if len(b) > 0 {
			if err = json.Unmarshal(b, userGroup); err != nil {
				return nil, err
			}

			return userGroup, nil
		}
	}

	//缓存不存在，从数据库查找
	if userGroup, err = u.model.FindByGid(gid); err != nil {
		return nil, err
	}

	//保存缓存
	if b, err = json.Marshal(userGroup); err == nil {
		core.CLusterClient.Set(redisKey, string(b), core.DefaultExpire)
	}

	return userGroup, nil
}

//修改群信息
func (u *UserGroupsService) EditGroup(res *EditGroupRequest) (err error) {
	var (
		userGroup = &model.UserGroups{
			Uuid:        res.Uuid,
			Gid:         res.Gid,
			Name:        res.Name,
			Description: res.Description,
			Avatar:      res.Avatar,
		}
		redisKey = "user_groups:gid:" + res.Gid + ":json"
	)
	if err = u.model.EditGroups(userGroup); err != nil {
		return err
	}

	//删除缓存
	core.CLusterClient.Del(redisKey)

	return err
}
