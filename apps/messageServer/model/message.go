package model

import (
	"github.com/xormplus/xorm"
	"time"
)

type (
	UserMessage struct {
		Id         int64
		Token      string    `json:"token" xorm:"varchar(125) notnull index"` //两个uuid经过加密得到的token
		FromUuid   string    `json:"from_uuid" xorm:"varchar(125) notnull"`   //发送消息的人的uuid
		ToUuid     string    `json:"to_uuid" xorm:"varchar(125) notnull"`     //接收消息的用户的uuid
		Type       string    `json:"type" xorm:"varchar(25) default('text')"` //消息类型 text、img、music、movie
		Message    string    `json:"message" xorm:"varchar(255) notnull"`     //消息内容
		CreateTime time.Time `json:"-" xorm:"created"`
	}
	UserMessageModel struct {
		engine *xorm.Engine
	}
)

func NewUserMessageModel(x *xorm.Engine) *UserMessageModel {
	return &UserMessageModel{
		engine: x,
	}
}
