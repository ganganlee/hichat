package service

import (
	"encoding/json"
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
	var (
		redisKey      = "historyRecord:uuid:" + uuid + ":hash"
		unread   uint = 0
		err      error
		b        []byte
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
		unread += 1
	}

	val.Unread = unread

	//将消息转为字符串
	if b, err = json.Marshal(val); err != nil {
		return
	}
	err = core.CLusterClient.HSet(redisKey, key, string(b)).Err()

	//判断当消息内容为空时，结束方法；不为空时，将消息加入之历史消息
	if val.Content == "" {
		return
	}

	//定义缓存key
	redisKey = "historyRecord:uuid:" + uuid + ":to_id:" + val.Id + ":list"
	core.CLusterClient.LPush(redisKey, string(b))
	//历史消息只保留100条数据
	core.CLusterClient.LTrim(redisKey, 0, 100)
}

//获取历史聊天内容
func (h *HistoryRecord) HistoryInfo(id string) {
	var (
		redisKey string
		list     []string
	)

	//参数判断
	if id == "" {
		core.ResponseSocketMessage(h.conn, "err", "用户id不能为空")
		return
	}

	//定义缓存key
	redisKey = "historyRecord:uuid:" + h.uuid + ":to_id:" + id + ":list"

	//获取缓存
	list = core.CLusterClient.LRange(redisKey, 0, -1).Val()
	core.ResponseSocketMessage(h.conn, "HistoryInfo", list)

	//缓存当前用户的聊天对象，用来作为增加未读消息的判断依据
	redisKey = "historyRecord:uuid:" + h.uuid + ":messageUser:string"
	core.CLusterClient.Set(redisKey, id, 5*time.Minute)
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

//清除未读消息
func (h *HistoryRecord) ClearUnread(id string) {
	var (
		redisKey = "historyRecord:uuid:" + h.uuid + ":hash"
		err      error
		b        []byte
		msg      *HistoryMessage
	)

	//参数判断
	if id == "" {
		return
	}

	//获取缓存
	if b, err = core.CLusterClient.HGet(redisKey, id).Bytes(); err != nil {
		return
	}

	//将字符串解析为结构体
	msg = new(HistoryMessage)
	if err = json.Unmarshal(b, msg); err != nil {
		return
	}

	//修改未读消息数量
	msg.Unread = 0

	//将消息转化为字符串保存
	if b, err = json.Marshal(msg); err != nil {
		return
	}

	core.CLusterClient.HSet(redisKey, id, string(b))
}
