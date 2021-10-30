package main

import (
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/fileService/common"
	"hichat.zozoo.net/apps/fileService/route"
	"hichat.zozoo.net/middleware"
	"log"
	"os"
)

//文件服务
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

	//连接数据库
	if _, err = common.OrmConn(cfg); err != nil {
		log.Fatalf("数据库连接失败 err:%v\n", err)
	}

	//启动框架，注册路由
	r := gin.Default()
	//设置静态文件访问地址
	r.Static("/static", "./static")

	//允许跨域
	r.Use(middleware.Cors(),middleware.Auth())

	//注册路由
	route.InitRoute(r)

	//启动框架
	if err = r.Run(cfg.Host); err != nil {
		log.Printf("框架启动失败 err:%v\n", err)
	}
}
