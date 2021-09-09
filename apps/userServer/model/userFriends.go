package model

import (
	"github.com/xormplus/xorm"
	"time"
)

//用户好友数据模型
type (
	UserFriends struct {
		Id         int64     `json:"id"`
		Uuid       string    `json:"uuid" xorm:"varchar(125) notnull unique(index::uuid:friend_uuid)"`
		FriendUuid string    `json:"friend_uuid" xorm:"varchar(255) notnull unique(index::uuid:friend_uuid)"`
		Status     int32     `json:"status" xorm:"tinyint(1) default(0)"` //好友状态 0:未添加，1:已添加
		CreateTime time.Time `json:"-" xorm:"created"`
		UpdateTime time.Time `json:"-" xorm:"updated"`
	}
	UserFriendUser struct {
		UserFriends `xorm:"extends"`
		User        `xorm:"extends"`
	}
	UserFriendsModel struct {
		engine *xorm.Engine
	}
)

func NewUserFriendsModel(e *xorm.Engine) *UserFriendsModel {
	return &UserFriendsModel{
		e,
	}
}

//申请添加好友
func (u *UserFriendsModel) ApplyFriends(friend *UserFriends) (err error) {
	_, err = u.engine.Insert(friend)
	return err
}

//同意好友申请
func (u *UserFriendsModel) ApproveFriends(uuid string, friendUuid string) (err error) {
	session := u.engine.NewSession()
	defer session.Close()

	//开启事务
	if err = session.Begin(); err != nil {
		return err
	}

	var friend = new(UserFriends)
	friend.Status = 1

	//1、修改之前申请的状态
	if _, err = u.engine.Where("uuid=? and friend_uuid=?", uuid, friendUuid).Update(friend); err != nil {
		session.Rollback()
		return err
	}

	//2、添加自己是对方好友的数据
	friend.Uuid = friendUuid
	friend.FriendUuid = uuid
	if _, err = u.engine.Insert(friend); err != nil {
		session.Rollback()
		return err
	}

	//提交事务
	return session.Commit()
}

//删除好友
func (u *UserFriendsModel) DelFriends(uuid string, friendUuid string) (err error) {
	session := u.engine.NewSession()

	var friend = new(UserFriends)

	//开启事务
	if err = session.Begin(); err != nil {
		return err
	}
	//1、将对方从自己的好友列表删除
	if _, err = u.engine.Where("uuid=? and friend_uuid=?", uuid, friendUuid).Delete(friend); err != nil {
		session.Rollback()
		return err
	}

	//2、将自己从对方列表中删除
	if _, err = u.engine.Where("uuid=? and friend_uuid=?",friendUuid,uuid).Delete(friend); err != nil {
		session.Rollback()
		return err
	}

	//提交事务
	return session.Commit()
}

//获取好友列表
func (u *UserFriendsModel) Friends(uuid string) (list []UserFriendUser, err error) {
	userFriends := make([]UserFriendUser, 0)
	err = u.engine.Table("user_friends").Join("INNER", "user", "user_friends.friend_uuid=user.uuid").Cols("user_friends.status", "user.uuid", "user.username", "user.avatar").Where("user_friends.uuid=?", uuid).Find(&userFriends)
	if err != nil {
		return nil, err
	}
	return userFriends, nil
}
