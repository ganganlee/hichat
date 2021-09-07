package common

import (
	"fmt"
	"github.com/micro/go-micro/v2/config"
)

type (
	Etcd struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	RpcServer struct {
		UserRpc string `json:"user_rpc"`
	}
	Config struct {
		Version    string    `json:"version"`
		ServerName string    `json:"server_name"`
		RpcServer  RpcServer `json:"rpc_server"`
		Host       string    `json:"host"`
		Etcd       Etcd      `json:"etcd"`
	}
)

var AppCfg *Config

//读取配置文件
func GetConfig(path string) (cfg *Config, err error) {

	//加载配置文件
	if err = config.LoadFile(path); err != nil {
		return nil, err
	}

	//将配置文件转换为结构体
	cfg = new(Config)
	if err = config.Get().Scan(cfg); err != nil {
		return nil, err
	}

	fmt.Println(cfg)

	AppCfg = cfg
	return cfg, nil
}
