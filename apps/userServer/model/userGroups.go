package model

import (
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
