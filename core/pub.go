package core

import (
	"encoding/json"
	"errors"
	"reflect"
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
		return errors.New("rpc消息传入错误")
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
