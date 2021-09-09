package rpc

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/apps/userServer/service"
	"hichat.zozoo.net/rpc/userFriends"
)

type UserFriendsRpc struct {
	service *service.UserFriendsService
}

func NewUserFriendsRpc(s *service.UserFriendsService) *UserFriendsRpc {
	return &UserFriendsRpc{
		s,
	}
}

//获取好友列表
func (u *UserFriendsRpc) Friends(ctx context.Context, res *userFriends.FriendsRequest, rsp *userFriends.FriendsResponse) error {

	var (
		list    []model.UserFriendUser
		friends []*userFriends.Friend
		err     error
	)

	if res.Uuid == "" {
		return errors.New("用户id不能为空")
	}

	//TODO 判断缓存

	//获取好友列表
	if list, err = u.service.Friends(res.Uuid); err != nil {
		return err
	}

	friends = make([]*userFriends.Friend, 0)

	//组织返回数据
	for _, friend := range list {
		friends = append(friends, &userFriends.Friend{
			Uuid:     friend.User.Uuid,
			Username: friend.Username,
			Avatar:   friend.Avatar,
			Status:   friend.Status,
		})
	}

	rsp.Friends = friends

	return nil
}

//添加好友申请
func (u *UserFriendsRpc) ApplyFriends(ctx context.Context, res *userFriends.ApplyFriendsRequest, rsp *userFriends.ApplyFriendsResponse) error {
	var (
		param = &service.ApplyFriendsRequest{
			Uuid:       res.Uuid,
			FriendUuid: res.FriendUuid,
		}
		err error
	)

	//验证参数
	validate := validator.New()
	if err = validate.Struct(param); err != nil {
		return err
	}

	//调用服务方法
	if err = u.service.ApplyFriends(param); err != nil {
		return err
	}

	//TODO 删除缓存
	rsp.Msg = "ok"
	return nil
}

//同意好友申请
func (u *UserFriendsRpc) ApproveFriends(ctx context.Context, res *userFriends.ApproveFriendsRequest, rsp *userFriends.ApproveFriendsResponse) error {
	var (
		param = &service.ApplyFriendsRequest{
			Uuid:       res.Uuid,
			FriendUuid: res.FriendUuid,
		}
		err error
	)

	//参数验证
	validate := validator.New()
	if err = validate.Struct(param); err != nil {
		return err
	}

	//调用同意方法
	if err = u.service.ApproveFriends(param); err != nil {
		return err
	}

	//TODO 删除缓存
	rsp.Msg = "ok"
	return nil
}

//删除好友
func (u *UserFriendsRpc) DelFriends(ctx context.Context, res *userFriends.DelFriendsRequest, rsp *userFriends.DelFriendsResponse) error {
	var (
		param = &service.ApplyFriendsRequest{
			Uuid:       res.Uuid,
			FriendUuid: res.FriendUuid,
		}
		err error
	)

	//参数验证
	validate := validator.New()
	if err = validate.Struct(param); err != nil {
		return err
	}

	//调用参数方法
	if err = u.service.DelFriends(param); err != nil {
		return err
	}

	//TODO 删除缓存
	rsp.Msg = "ok"
	return nil
}
