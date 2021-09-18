package model

import (
	"github.com/xormplus/xorm"
	"time"
)

//用户群成员数据模型
type (
	UserGroupMembers struct {
		Id         int64     `json:"id"`
		GroupId    string    `json:"group_id" xorm:"notnull index(index:group_gid:user_uuid) unique(unique:group_gid:user_uuid)"`
		UserId     string    `json:"user_id" xorm:"notnull index(index:group_gid:user_uuid) unique(unique:group_gid:user_uuid)"`
		CreateTime time.Time `json:"-" xorm:"created"`
		UpdateTime time.Time `json:"-" xorm:"updated"`
	}
	UserGroupMembersModel struct {
		engine *xorm.Engine
	}
)

func NewUserGroupMembersModel(e *xorm.Engine) *UserGroupMembersModel {
	return &UserGroupMembersModel{
		e,
	}
}

//添加群成员
func (u *UserGroupMembersModel) InsertMember(member *UserGroupMembers) (err error) {
	_, err = u.engine.Insert(member)
	return err
}

//删除群成员
func (u *UserGroupMembersModel) DelByMemberId(member *UserGroupMembers) (err error) {
	_, err = u.engine.Delete(member)
	return err
}

//删除所有群成员
func (u *UserGroupMembersModel) DelByGroupId(groupId string) (err error) {
	var m = new(UserGroupMembers)
	_, err = u.engine.Where("group_id=?", groupId).Delete(m)
	return err
}

//获取成员群列表
func (u *UserGroupMembersModel) MemberGroups(uuid string) (list []UserGroupMembers, err error) {
	list = make([]UserGroupMembers, 0)
	if err = u.engine.Where("user_id=?", uuid).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

//获取群成员列表
func (u *UserGroupMembersModel) Members(gid string) (list []UserGroupMembers, err error) {
	list = make([]UserGroupMembers, 0)
	if err = u.engine.Where("group_id=?", gid).Cols("user_id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}
