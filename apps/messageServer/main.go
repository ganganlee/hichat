package main

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/messageServer/common"
	"hichat.zozoo.net/apps/messageServer/route"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/middleware"
	"log"
	"os"
)

//消息服务

func main() {
	var (
		path string
		err  error
		cfg  *common.Config
	)

	//获取当前目录
	if path, err = os.Getwd(); err != nil {
		log.Fatalf("获取配置文件目录失败 err:%v\n", err)
	}
	path = path + "/config/app.json"

	//获取配置文件
	if cfg, err = common.GetConfig(path); err != nil {
		log.Fatalf("读取配置文件失败 err:%v\n", err)
	}

	//连接redis
	core.RedisClusterConn(cfg.Redis.Address, cfg.Redis.MaxActive, cfg.Redis.MinIdle, cfg.Redis.Password)

	//启动框架，注册路由
	r := gin.Default()

	//允许跨域
	r.Use(middleware.Cors())

	//注册路由
	route.InitListenRoute(r)

	//启动框架
	if err = r.Run(cfg.Host); err != nil {
		log.Printf("框架启动失败 err:%v\n", err)
	}
}
