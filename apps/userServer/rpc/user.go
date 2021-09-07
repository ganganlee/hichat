package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/apps/userServer/service"
	"hichat.zozoo.net/rpc/user"
)

type (
	UserRpc struct {
		service *service.UserService
	}

	//登录结构体
	LoginRequest struct {
		Username string `json:"username" validate:"required,min=3,max=25"`
		Password string `json:"password" validate:"required,min=6,max=25"`
	}
)

func NewUserRpc(s *service.UserService) *UserRpc {
	return &UserRpc{
		s,
	}
}

//用户注册
func (u *UserRpc) Register(ctx context.Context, res *user.RegisterRequest, rsp *user.RegisterResponse) error {
	var (
		err  error
		user *model.User
		uuid string
	)

	//实例化用户模型
	user = new(model.User)

	if res.User == nil {
		return errors.New("参数错误")
	}

	user.Username = res.User.Username
	user.Password = res.User.Password
	user.Avatar = res.User.Avatar

	//参数验证
	validate := validator.New()
	if err = validate.Struct(user); err != nil {
		return err
	}

	//提交注册
	if uuid, err = u.service.Register(user); err != nil {
		return err
	}
	rsp.Uuid = uuid
	return nil
}

//用户登录
func (u *UserRpc) Login(ctx context.Context, res *user.LoginRequest, rsp *user.LoginResponse) error {
	var (
		userOrm      *model.User
		loginRequest *LoginRequest
		err          error
		exist        bool
	)

	loginRequest = new(LoginRequest)
	loginRequest.Username = res.Username
	loginRequest.Password = res.Password

	//参数验证
	validate := validator.New()
	if err = validate.Struct(loginRequest); err != nil {
		return err
	}

	//实例化数据模型
	userOrm = new(model.User)
	userOrm.Username = loginRequest.Username
	userOrm.Password = loginRequest.Password

	if exist, err = u.service.Login(userOrm); err != nil {
		return err
	}

	if !exist {
		return errors.New("用户名或密码错误")
	}

	rsp.User = new(user.User)
	rsp.User.Uuid = userOrm.Uuid

	return nil
}

//根据用户uuid查找用户
func (u *UserRpc) FindByUuid(ctx context.Context, res *user.FindByUuidRequest, rsp *user.FindByUuidResponse) error {
	var (
		err     error
		rpcUser *model.User
	)
	if rpcUser, err = u.service.FindByUuid(res.Uuid); err != nil {
		return err
	}

	//定义响应数据
	rsp.User = &user.User{
		Uuid:     rpcUser.Uuid,
		Username: rpcUser.Username,
		Avatar:   rpcUser.Avatar,
	}

	return nil
}

//修改用户信息
func (u *UserRpc) EditInfo(ctx context.Context, res *user.EditInfoRequest, rsp *user.EditInfoResponse) error {
	var (
		editRequest *service.EditRequest
		err         error
	)

	editRequest = &service.EditRequest{
		Uuid:     res.User.Uuid,
		Username: res.User.Username,
		Password: res.User.Password,
		Avatar:   res.User.Avatar,
	}

	//参数验证
	validate := validator.New()
	if err = validate.Struct(editRequest); err != nil {
		return err
	}

	//调用方法保存操作
	if err = u.service.Edit(editRequest); err != nil {
		return err
	}

	rsp.Msg = "ok"

	return nil
}

//根据用户名查找用户
func (u *UserRpc) FindByUsername(ctx context.Context, res *user.FindByUsernameRequest, rsp *user.FindByUsernameResponse) error {

	var (
		users []model.User
		list  []*user.User
		err   error
	)

	if res.Username == "" {
		return errors.New("请输入用户名进行查找")
	}

	if users, err = u.service.FindByUsername(res.Username); err != nil {
		return err
	}

	fmt.Println(users)
	list = make([]*user.User, 0)
	for _, val := range users {
		fmt.Println(val)
		list = append(list, &user.User{
			Uuid:     val.Uuid,
			Username: val.Username,
			Avatar:   val.Avatar,
		})
	}

	rsp.Users = list
	return nil
}
