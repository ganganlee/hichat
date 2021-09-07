package main

//网站访问入口

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/webServer/common"
	"hichat.zozoo.net/apps/webServer/route"
	"log"
	"os"
)

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

	//启动框架，注册路由
	r := gin.Default()
	route.InitWebRoute(r)

	//启动框架
	if err = r.Run(cfg.Host); err != nil {
		log.Printf("框架启动失败 err:%v\n", err)
	}
}
