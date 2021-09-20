package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"hichat.zozoo.net/core"
	"time"
)

//好友聊天列表记录
type (
	HistoryRecord struct {
		conn *websocket.Conn
		uuid string
	}

	//客户端添加聊天历史列表
	SetHistoryRecordRequest struct {
		Type string `json:"type" validate:"required"`
		Id   string `json:"id" validate:"required"`
	}

	//消息内容结构体
	HistoryMessage struct {
		Id          string `json:"id"`           //聊天对象id
		MessageType string `json:"message_type"` //聊天消息类型 privateMessage:私聊，groupMessage:群聊消息，noticeMessage:通知消息
		ContentType string `json:"content_type"` //消息内容类型 text,voice,video,img,icon
		Content     string `json:"content"`      //消息内容
		Date        string `json:"date"`         //时间
		Unread      uint   `json:"unread"`       //消息未读数量
		Name        string `json:"name"`
		Avatar      string `json:"avatar"`
		Uuid        string `json:"uuid"`
	}
)

func NewHistoryRecord(conn *websocket.Conn, uuid string) *HistoryRecord {
	return &HistoryRecord{
		conn: conn,
		uuid: uuid,
	}
}

//获取聊天记录列表
func (h *HistoryRecord) List(str string) {
	var (
		redisKey string
		mp       map[string]string
	)

	mp = make(map[string]string, 0)
	redisKey = "historyRecord:uuid:" + h.uuid + ":hash"
	mp = core.CLusterClient.HGetAll(redisKey).Val()
	core.ResponseSocketMessage(h.conn, "HistoryRecord", mp)
}

//设置聊天列表记录
func (h *HistoryRecord) SetHistoryRecord(content string) {
	var (
		res     = new(SetHistoryRecordRequest)
		err     error
		history *HistoryMessage
		t       = time.Now()
	)

	if err = json.Unmarshal([]byte(content), res); err != nil {
		core.ResponseSocketMessage(h.conn, "err", err.Error())
		return
	}

	history = &HistoryMessage{
		Id:          res.Id,
		MessageType: res.Type,
		ContentType: "text",
		Uuid:        h.uuid,
		Content:     "",
		Date:        t.Format("2006-01-02 15:04:05"),
	}

	//将消息加入到缓存
	h.PushHistoryRecord(h.uuid, res.Id, history, false)
}

//将消息加入到聊天列表缓存中
func (h *HistoryRecord) PushHistoryRecord(uuid string, key string, val *HistoryMessage, addUnread bool) {
	var redisKey = "historyRecord:uuid:" + uuid + ":hash"

	var (
		unread uint = 0
		err    error
		b      []byte
	)

	//判断需要增加未读消息
	if addUnread {
		//获取缓存，判断是否存在，缓存存在时给缓存未读消息数量加一
		var (
			str     string
			history *HistoryMessage
		)
		history = new(HistoryMessage)
		str = core.CLusterClient.HGet(redisKey, key).Val()
		if str != "" {
			if err = json.Unmarshal([]byte(str), history); err == nil {
				unread = history.Unread
			}
		}

	}

	val.Unread = unread

	//将消息转为字符串
	if b, err = json.Marshal(val); err != nil {
		return
	}

	err = core.CLusterClient.HSet(redisKey, key, string(b)).Err()
	fmt.Println(err)
}

//根据id从聊天列表中删除消息
func (h *HistoryRecord) RemoveHistoryRecord(id string) {
	var (
		err      error
		redisKey string
	)

	if id == "" {
		core.ResponseSocketMessage(h.conn, "err", "消息id不能为空")
		return
	}

	redisKey = "historyRecord:uuid:" + h.uuid + ":hash"
	if err = core.CLusterClient.HDel(redisKey, id).Err(); err != nil {
		core.ResponseSocketMessage(h.conn, "err", err.Error())
		return
	}
}
