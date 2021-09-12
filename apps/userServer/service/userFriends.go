package service

import (
	"encoding/json"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/core"
)

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
	var (
		redisKey = "userFriends:uuid:" + uuid + ":list"
		b        []byte
	)

	//判断缓存是否存在
	if b, err = core.CLusterClient.Get(redisKey).Bytes(); err == nil {
		if len(b) > 0 {
			//缓存存在，将缓存解析为对象返回
			list = make([]model.UserFriendUser, 0)
			if err = json.Unmarshal(b, &list); err != nil {
				return nil, err
			}

			return list, nil
		}
	}

	//缓存不存在，调用sql查询数据
	if list, err = u.engine.Friends(uuid); err != nil {
		return nil, err
	}

	//判断数据存在，则将数据加入到缓存中
	if len(list) == 0 {
		return list, nil
	}

	if b, err = json.Marshal(list); err != nil {
		return nil, err
	}

	core.CLusterClient.Set(redisKey, string(b), core.DefaultExpire)

	return list, nil
}

//申请添加好友
func (u *UserFriendsService) ApplyFriends(res *ApplyFriendsRequest) (err error) {
	//自己添加对方为好友，所以是在对方数据中添加自己为好友，这样对方在查看好友时就能看到你的申请
	var (
		friend = &model.UserFriends{
			Uuid:       res.FriendUuid,
			FriendUuid: res.Uuid,
		}
		redisKey string
	)
	if err = u.engine.ApplyFriends(friend); err != nil {
		return err
	}

	//清除对方的缓存，这样对方在获取列表时就能看见好友申请
	redisKey = "userFriends:uuid:" + res.FriendUuid + ":list"
	core.CLusterClient.Del(redisKey)

	return nil
}

//同意添加好友
func (u *UserFriendsService) ApproveFriends(res *ApplyFriendsRequest) (err error) {

	if err = u.engine.ApproveFriends(res.Uuid, res.FriendUuid); err != nil {
		return err
	}

	//清除双方好友列表缓存
	var redisKey = "userFriends:uuid:" + res.FriendUuid + ":list"
	core.CLusterClient.Del(redisKey)
	redisKey = "userFriends:uuid:" + res.Uuid + ":list"
	core.CLusterClient.Del(redisKey)

	return nil
}

//拒绝好友申请
func (u *UserFriendsService) RefuseFriend(res *ApplyFriendsRequest) (err error) {
	if err = u.engine.RefuseFriend(res.Uuid, res.FriendUuid); err != nil {
		return err
	}

	//删除缓存
	var redisKey = "userFriends:uuid:" + res.Uuid + ":list"
	core.CLusterClient.Del(redisKey)
	return nil
}

//删除好友
func (u *UserFriendsService) DelFriends(res *ApplyFriendsRequest) (err error) {
	if err = u.engine.DelFriends(res.Uuid, res.FriendUuid); err != nil {
		return err
	}

	//清除双方好友列表缓存
	var redisKey = "userFriends:uuid:" + res.FriendUuid + ":list"
	core.CLusterClient.Del(redisKey)
	redisKey = "userFriends:uuid:" + res.Uuid + ":list"
	core.CLusterClient.Del(redisKey)

	return nil
}
