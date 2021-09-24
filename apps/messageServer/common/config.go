package common

import "github.com/micro/go-micro/v2/config"

type (
	Etcd struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	Mysql struct {
		Host      string `json:"host"`
		Database  string `json:"database"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Charset   string `json:"charset"`
		ShowSql   bool   `json:"show_sql"`
		MaxIdle   int    `json:"max_idle"`
		MaxActive int    `json:"max_active"`
	}
	Redis struct {
		MinIdle   int      `json:"min_idle"`
		MaxActive int      `json:"max_active"`
		Password  string   `json:"password"`
		Address   []string `json:"address"`
	}
	RpcServer struct {
		UserRpc    string `json:"user_rpc"`
		GatewayRpc string `json:"gateway_rpc"`
	}
	Config struct {
		Version    string    `json:"version"`
		ServerName string    `json:"server_name"`
		Host       string    `json:"host"`
		MqHost     string    `json:"mq_host"`
		RpcServer  RpcServer `json:"rpc_server"`
		Etcd       Etcd      `json:"etcd"`
		Mysql      Mysql     `json:"mysql"`
		Redis      Redis     `json:"redis"`
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

	AppCfg = cfg
	return cfg, nil
}
