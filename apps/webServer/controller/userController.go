package controller

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/webServer/common"
	"hichat.zozoo.net/apps/webServer/service"
	"hichat.zozoo.net/core"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{
		s,
	}
}

//用户注册
func (u *UserController) Register(c *gin.Context) {
	var (
		err error
		res *service.RegisterRequest
		rsp *service.RegisterResponse
	)

	res = new(service.RegisterRequest)

	if err = c.ShouldBindJSON(res); err != nil {
		core.ResponseError(c, common.Translate(err))
		return
	}

	if rsp, err = u.service.Register(res); err != nil {
		core.ResponseError(c, err.Error())
		return
	}

	core.ResponseSuccess(c, rsp)
}

//用户登录
func (u *UserController) Login(c *gin.Context) {
	var (
		err error
		res *service.LoginRequest
		rsp *service.LoginResponse
	)

	res = new(service.LoginRequest)

	//参数验证
	if err = c.ShouldBindJSON(res); err != nil {
		core.ResponseError(c, common.Translate(err))
		return
	}

	//提交登录
	if rsp, err = u.service.Login(res); err != nil {
		core.ResponseError(c, err.Error())
		return
	}

	core.ResponseSuccess(c, rsp)
}

//根据token查找用户
func (u *UserController) FindByUuid(c *gin.Context) {
	var (
		err error
		rsp *service.FindByUuid
	)

	if rsp, err = u.service.FindByUuid(); err != nil {
		core.ResponseError(c, err.Error())
		return
	}

	core.ResponseSuccess(c, rsp)
}
