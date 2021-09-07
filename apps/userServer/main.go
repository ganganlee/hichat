package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	xorm2 "github.com/xormplus/xorm"
	"hichat.zozoo.net/apps/userServer/common"
	"hichat.zozoo.net/apps/userServer/model"
	"hichat.zozoo.net/apps/userServer/rpc"
	service2 "hichat.zozoo.net/apps/userServer/service"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/user"
	"log"
	"os"
)

//用户rpc服务
func main() {
	var (
		path    string
		err     error
		cfg     *common.Config
		userRpc *rpc.UserRpc
		service micro.Service
		xorm    *xorm2.Engine
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
	if xorm, err = common.OrmConn(cfg); err != nil {
		log.Fatalf("数据库连接失败 err:%v\n", err)
	}

	//连接redis
	core.RedisClusterConn(cfg.Redis.Address, cfg.Redis.MaxActive, cfg.Redis.MinIdle, cfg.Redis.Password)

	//初始化微服务
	service = micro.NewService(
		micro.Name(cfg.ServerName),
		micro.Registry(etcd.NewRegistry(registry.Addrs(cfg.Etcd.Host))),
	)
	service.Init()

	//注册用户服务
	userModel := model.NewUserModel(xorm)
	userSvc := service2.NewUserService(userModel)
	userRpc = rpc.NewUserRpc(userSvc)
	if err = user.RegisterUserServiceHandler(service.Server(), userRpc); err != nil {
		log.Fatalf("注册用户服务失败 err:%v\n", err)
	}

	//运行微服务
	if err = service.Run(); err != nil {
		log.Fatalf("微服务启动失败 err:%v\n", err)
	}
}

//func init() {
//	// 获取日志文件句柄
//	// 已 只写入文件|没有时创建|文件尾部追加 的形式打开这个文件
//	logFile, err := os.OpenFile(`./debug.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		panic(err)
//	}
//
//	// 设置存储位置
//	log.SetOutput(logFile)
//}
