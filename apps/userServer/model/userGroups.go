package model

import (
	"errors"
	"github.com/xormplus/xorm"
	"time"
)

//用户群数据模型
type (
	UserGroups struct {
		Id          int64
		Gid         string    `json:"gid" xorm:"varchar(125) unique unique(index:gid:uuid) notnull"`
		Uuid        string    `json:"uuid" xorm:"varchar(125) notnull unique(index:gid:uuid)"`
		Name        string    `json:"name" xorm:"varchar(125) notnull"`
		Description string    `json:"description" xorm:"default('')"`
		Avatar      string    `json:"avatar" xorm:"varchar(125) notnull"`
		CreateTime  time.Time `json:"-" xorm:"created"`
		UpdateTime  time.Time `json:"-" xorm:"updated"`
	}
	UserGroupsModel struct {
		engine *xorm.Engine
	}
)

func NewUserGroupsModel(e *xorm.Engine) *UserGroupsModel {
	return &UserGroupsModel{
		e,
	}
}

//创建群
func (u *UserGroupsModel) Create(g *UserGroups) (err error) {
	_, err = u.engine.Insert(g)
	return err
}

//删除群
func (u *UserGroupsModel) DelByModel(g *UserGroups) (err error) {
	_, err = u.engine.Delete(g)
	return err
}

//根据gid查找群
func (u *UserGroupsModel) FindByGid(gid string) (userGroups *UserGroups, err error) {
	var exist bool
	userGroups = new(UserGroups)

	if exist, err = u.engine.Where("gid=?", gid).Get(userGroups); err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("群不存在")
	}

	return
}

//修改群信息
func (u *UserGroupsModel) EditGroups(userGroup *UserGroups) (err error) {
	_, err = u.engine.Where("gid=? AND uuid=?", userGroup.Gid, userGroup.Uuid).Cols("name", "description", "avatar").Update(userGroup)
	return err
}
