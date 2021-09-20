package common

import "github.com/micro/go-micro/v2/config"

type (
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
	Config struct {
		Version    string `json:"version"`
		ServerName string `json:"server_name"`
		Host       string `json:"host"`
		Mysql      Mysql  `json:"mysql"`
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
