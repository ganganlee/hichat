package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	limiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	xorm2 "github.com/xormplus/xorm"
	"hichat.zozoo.net/apps/gatewayServer/common"
	"hichat.zozoo.net/apps/gatewayServer/models"
	"hichat.zozoo.net/apps/gatewayServer/rpc"
	service2 "hichat.zozoo.net/apps/gatewayServer/service"
	"hichat.zozoo.net/core"
	gateway "hichat.zozoo.net/rpc/Gateway"
	"log"
	"os"
)

//用户rpc服务
func main() {
	var (
		path       string
		err        error
		cfg        *common.Config
		gatewayRpc *rpc.GatewayRpc
		service    micro.Service
		xorm       *xorm2.Engine
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

	//连接mq服务器
	common.InitRabbitMq(cfg)

	//初始化微服务
	service = micro.NewService(
		micro.Name(cfg.ServerName),
		micro.Registry(etcd.NewRegistry(registry.Addrs(cfg.Etcd.Host))),
		//限流
		micro.WrapHandler(
			limiter.NewHandlerWrapper(10),
			),
	)
	service.Init()

	//注册网关服务
	messageModel := models.NewMessageModel(xorm)
	gatewayService := service2.NewGatewayService(messageModel)
	gatewayRpc = rpc.NewGatewayRpc(gatewayService)
	if err = gateway.RegisterGatewayServiceHandler(service.Server(), gatewayRpc); err != nil {
		log.Fatalf("注册消息服务失败 err:%v\n", err)
	}

	//运行微服务
	if err = service.Run(); err != nil {
		log.Fatalf("微服务启动失败 err:%v\n", err)
	}
}
