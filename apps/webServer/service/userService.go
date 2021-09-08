package service

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/apps/webServer/common"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/user"
	"math/rand"
)

type (
	UserService struct {
		userRpc user.UserService
	}

	//注册请求结构体
	RegisterRequest struct {
		Username string `json:"username" binding:"required,min=3,max=25"`
		Password string `json:"password" binding:"required,min=6,max=25"`
		Avatar   string `json:"avatar" binding:"required,url"`
	}
	RegisterResponse struct {
		Uuid string `json:"uuid"`
	}

	//登录请求结构体
	LoginRequest struct {
		Username string `json:"username" binding:"required,min=3,max=25"`
		Password string `json:"password" binding:"required,min=6,max=25"`
	}
	LoginResponse struct {
		Token string `json:"token"`
	}

	FindByUuid struct {
		Uuid     string
		Username string
		Avatar   string
		Host     string
	}
)

func NreUserService() *UserService {
	//注册用户rpc服务
	service := micro.NewService(
		micro.Registry(etcd.NewRegistry(registry.Addrs(common.AppCfg.Etcd.Host))),
	)
	userRpc := user.NewUserService(common.AppCfg.RpcServer.UserRpc, service.Client())
	return &UserService{
		userRpc: userRpc,
	}
}

//用户注册
func (u *UserService) Register(res *RegisterRequest) (rsp *RegisterResponse, err error) {

	var (
		rpcRsp  *user.RegisterResponse
		rpcUser = &user.User{
			Username: res.Username,
			Password: res.Password,
			Avatar:   res.Avatar,
		}

		rpcRes = new(user.RegisterRequest)
	)
	rpcRes.User = rpcUser

	//提交到rpc注册用户
	if rpcRsp, err = u.userRpc.Register(context.TODO(), rpcRes); err != nil {
		return nil, core.DecodeRpcErr(err.Error())
	}

	rsp = new(RegisterResponse)
	rsp.Uuid = rpcRsp.Uuid

	return
}

//用户登录
func (u *UserService) Login(res *LoginRequest) (rsp *LoginResponse, err error) {
	var (
		rpcRes = &user.LoginRequest{
			Username: res.Username,
			Password: res.Password,
		}
		rpcRsp *user.LoginResponse
		token  string
	)

	//调用rpc方法
	if rpcRsp, err = u.userRpc.Login(context.TODO(), rpcRes); err != nil {
		return nil, core.DecodeRpcErr(err.Error())
	}

	//生成JWTToken
	if token, err = core.GenerateToken(rpcRsp.User.Uuid); err != nil {
		return nil, err
	}

	rsp = &LoginResponse{
		Token: token,
	}

	return
}

//根据uuid查找用户
func (u *UserService) FindByUuid(uuid string) (rsp *FindByUuid, err error) {
	var rpcRsp *user.FindByUuidResponse
	if rpcRsp, err = u.userRpc.FindByUuid(context.TODO(), &user.FindByUuidRequest{
		Uuid: uuid,
	}); err != nil {
		return nil, core.DecodeRpcErr(err.Error())
	}

	var messageServer = common.AppCfg.MessageHost[rand.Intn(len(common.AppCfg.MessageHost))]
	rsp = &FindByUuid{
		Uuid:     rpcRsp.User.Uuid,
		Username: rpcRsp.User.Username,
		Avatar:   rpcRsp.User.Avatar,
		Host:     messageServer,
	}
	return rsp, nil
}
