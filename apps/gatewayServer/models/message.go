package models

import (
	"github.com/xormplus/xorm"
	"time"
)

type (
	Message struct {
		Id          int64
		FromId      string    `json:"from_id" xorm:"varchar(125) notnull index index(index:from_id:to_id)"` //发送id
		ToId        string    `json:"to_id" xorm:"varchar(125) notnull index(index:from_id:to_id)"`         //接收id
		MsgType     string    `json:"msg_type" xorm:"varchar(25) default('text')"`                          //聊天类型
		ContentType string    `json:"content_type" xorm:"varchar(25) notnull"`                              //消息类型 text、img、music、movie
		Content     string    `json:"content" xorm:"varchar(255) notnull"`                                  //消息内容
		CreateTime  time.Time `json:"-" xorm:"created"`
	}
	MessageModel struct {
		engine *xorm.Engine
	}
)

func NewMessageModel(e *xorm.Engine) *MessageModel {
	return &MessageModel{
		e,
	}
}

//添加数据
func (m *MessageModel) Create(tableName string, msg *Message) (err error) {
	_, err = m.engine.Table(tableName).Insert(msg)
	return err
}
