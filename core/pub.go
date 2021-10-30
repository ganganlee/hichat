package core

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type RpcErr struct {
	Id     string `json:"id"`
	Code   int    `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

//解析rpc error
func DecodeRpcErr(msg string) (err error) {
	var rpcErr = new(RpcErr)
	if err := json.Unmarshal([]byte(msg), rpcErr); err != nil {
		return errors.New(msg)
	}

	return errors.New(rpcErr.Detail)
}

//通过方法名调用结构体方法
func CallFuncByName(obj interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	objVal := reflect.ValueOf(obj)
	m := objVal.MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), errors.New("方法不存在")
	}

	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	out = m.Call(in)
	return
}

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