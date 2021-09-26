package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/core"
)

type (
	UserService struct {
		model *model.UserModel
	}

	//修改信息
	EditRequest struct {
		Uuid     string `json:"uuid" validate:"required"`
		Username string `json:"username" validate:"min=3,max=25"`
		Password string `json:"password" validate:"min=6,max=25"`
		Avatar   string `json:"avatar" validate:"url"`
	}
)

func NewUserService(m *model.UserModel) *UserService {
	return &UserService{
		m,
	}
}

//用户注册
func (u *UserService) Register(user *model.User) (uid string, err error) {

	//创建用户uuid
	user.Uuid = uuid.New().String()

	//加密用户密码
	user.Password = fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))

	if err = u.model.Create(user); err != nil {
		return "", err
	}

	return user.Uuid, nil
}

//用户登录
func (u *UserService) Login(user *model.User) (exist bool, err error) {
	//加密用户密码
	user.Password = fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))
	return u.model.Login(user)
}

//根据用户uuid查找用户
func (u *UserService) FindByUuid(uuid string) (user *model.User, err error) {
	var (
		redisKey string
		userByte []byte
	)

	user = new(model.User)
	redisKey = "hichat_user:uuid:" + uuid + ":json"

	//判断缓存是否存在
	if userByte, err = core.CLusterClient.Get(redisKey).Bytes(); err == nil {
		//解析json
		if err = json.Unmarshal(userByte, user); err != nil {
			return nil, err
		}

		return user, nil
	}

	//查找数据库
	if user, err = u.model.FindByUuid(uuid); err != nil {
		return nil, err
	}

	//TODO 保存缓存
	if userByte, err = json.Marshal(user); err == nil {
		core.CLusterClient.Set(redisKey, string(userByte), core.DefaultExpire)
	}

	return user, nil
}

//修改用户信息
func (u *UserService) Edit(editRequest *EditRequest) (err error) {
	var (
		user     *model.User
		redisKey string
	)

	//赋值给结构体
	user = &model.User{
		Username: editRequest.Username,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(editRequest.Password))),
		Avatar:   editRequest.Avatar,
	}

	if err = u.model.Edit(user, editRequest.Uuid); err != nil {
		return err
	}

	//修改成功，删除之前的缓存
	redisKey = "hichat_user:uuid:" + editRequest.Uuid + ":json"
	core.CLusterClient.Del(redisKey)

	return nil
}

//根据用户名查找用户
func (u *UserService) FindByUsername(username string) (users []model.User, err error) {
	return u.model.FindByUsername(username)
}
