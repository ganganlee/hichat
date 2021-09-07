package core

import (
	"encoding/json"
	"errors"
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
