package common

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

//获取message table名称
func GetMessageTable(fromId string, toId string, msgType string) string {
	var (
		key string
	)

	//默认使用fromId加密作为数据表名称，当前消息类型为群聊时，使用fromId最为数据表名称
	key = fmt.Sprintf("%x", md5.Sum([]byte(toId)))

	//当消息类型不是群聊时，并且fromId的第一个字符大于toId的第一个字符，使用toId加密作为数据表名称
	if msgType != "groupMessage" {
		s1, s2 := fromId[0], toId[0]
		if s1 < s2 {
			key = fmt.Sprintf("%x", md5.Sum([]byte(toId)))
		} else {
			key = fmt.Sprintf("%x", md5.Sum([]byte(fromId)))
		}
	}

	key = string(key[0])
	if _, err := strconv.Atoi(key); err == nil {
		return "user_message" + key
	}
	return "user_message_" + key
}
