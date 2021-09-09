package service

import "hichat.zozoo.net/apps/userServer/model"

//好友服务方法
type (
	UserFriendsService struct {
		engine *model.UserFriendsModel
	}

	//添加好友申请
	ApplyFriendsRequest struct {
		Uuid       string `json:"uuid" validate:"required"`
		FriendUuid string `json:"friend_uuid" validate:"required"`
	}
)

func NewUserFriendsService(m *model.UserFriendsModel) *UserFriendsService {
	return &UserFriendsService{
		m,
	}
}

//获取好友列表
func (u *UserFriendsService) Friends(uuid string) (list []model.UserFriendUser, err error) {
	return u.engine.Friends(uuid)
}

//申请添加好友
func (u *UserFriendsService) ApplyFriends(res *ApplyFriendsRequest) (err error) {
	//自己添加对方为好友，所以是在对方数据中添加自己为好友，这样对方在查看好友时就能看到你的申请
	var friend = &model.UserFriends{
		Uuid:       res.FriendUuid,
		FriendUuid: res.Uuid,
	}
	return u.engine.ApplyFriends(friend)
}

//同意添加好友
func (u *UserFriendsService) ApproveFriends(res *ApplyFriendsRequest) (err error) {
	return u.engine.ApproveFriends(res.Uuid, res.FriendUuid)
}

//删除好友
func (u *UserFriendsService) DelFriends(res *ApplyFriendsRequest) (err error) {
	return u.engine.DelFriends(res.Uuid, res.FriendUuid)
}
