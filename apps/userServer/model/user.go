package model

import (
	"errors"
	"github.com/xormplus/xorm"
	"time"
)

type (
	User struct {
		Id       int64     `json:"id"`
		Uuid     string    `json:"uuid" xorm:"varchar(125) notnull unique"`
		Username string    `json:"username" xorm:"varchar(125) notnull unique" validate:"required,min=3,max=25"`
		Password string    `json:"password" xorm:"varchar(125) notnull" validate:"required,min=6,max=25"`
		Avatar   string    `json:"avatar" xorm:"notnull" validate:"required,url"`
		Created  time.Time `json:"-" xorm:"created"`
		Updated  time.Time `json:"-" xorm:"created"`
	}

	UserModel struct {
		engine *xorm.Engine
	}
)

func NewUserModel(e *xorm.Engine) *UserModel {
	return &UserModel{
		e,
	}
}

//创建用户
func (u *UserModel) Create(user *User) (err error) {
	_, err = u.engine.Insert(user)
	return err
}

//用户登录
func (u *UserModel) Login(user *User) (exist bool, err error) {
	return u.engine.Get(user)
}

//根据用户uuid查找用户
func (u *UserModel) FindByUuid(uuid string) (user *User, err error) {
	var (
		exist bool
	)

	user = new(User)

	if exist, err = u.engine.Where("uuid=?", uuid).Get(user); err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("用户不存在")
	}

	return
}

//修改用户信息
func (u *UserModel) Edit(user *User, uuid string) (err error) {
	_, err = u.engine.Where("uuid=?", uuid).Update(user)
	return err
}

//根据用户名查找用户
func (u *UserModel) FindByUsername(username string) (users []User, err error) {

	users = make([]User, 0)
	if err = u.engine.Where("username like ?", "%"+username+"%").Find(&users); err != nil {
		return nil, err
	}

	return
}
