package service

import (
	"github.com/google/uuid"
	"hichat.zozoo.net/apps/userServer/model"
)

type (
	UserGroupsService struct {
		model *model.UserGroupsModel
	}

	//获取群列表结构体
	GroupsRequest struct {
		Gid  string `json:"gid" validate:"required"`
		Uuid string `json:"uuid" validate:"required"`
	}
)

func NewUserGroupsService(m *model.UserGroupsModel) *UserGroupsService {
	return &UserGroupsService{
		m,
	}
}

//创建群
func (u *UserGroupsService) CreateGroup(g *model.UserGroups) (err error) {
	//创建用户gid
	g.Gid = uuid.New().String()
	return u.model.Create(g)
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
func (u *UserGroupsService) Groups(res *GroupsRequest) (err error) {
	return err
}
