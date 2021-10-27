package common

import "github.com/micro/go-micro/v2/config"

type (
	Etcd struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	Es struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Index    string `json:"index"`
	}
	CanalConfig struct {
		Address     string `json:"address"`
		Port        int    `json:"port"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Destination string `json:"destination"`
		SoTimeOut   int32  `json:"so_time_out"`
		IdleTimeOut int32  `json:"idle_time_out"`
	}
	Config struct {
		Version     string      `json:"version"`
		ServerName  string      `json:"server_name"`
		Etcd        Etcd        `json:"etcd"`
		Es          Es          `json:"es"`
		CanalConfig CanalConfig `json:"canal_config"`
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

//获取es创建索引的规则
func GetEsIndex() string {
	return `
		{
		  "settings": {
			"number_of_replicas": 0,
			"number_of_shards": 3,
			"analysis": {
			  "analyzer": {
				"ik":{
				  "tokenizer":"ik_max_word"
				}
			  }
			}
		  },
		  "mappings": {
			"properties": {
			  "id":{
				"type": "integer"
			  },
			  "from_id":{
				"type": "keyword"
			  },
			  "to_id":{
				"type": "keyword"
			  },
			  "msg_type":{
				"type": "keyword"
			  },
			  "content_type":{
				"type": "keyword"
			  },
			  "content":{
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_smart"
			  },
			  "create_time":{
				"type": "date"
			  }
			}
		  }
		}
	`
}
